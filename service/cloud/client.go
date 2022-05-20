package cloud

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Totally not stolen from https://github.com/couchbasecloud/rest-api-examples/blob/main/go/client/client.go

const (
	headerKeyTimestamp     = "Couchbase-Timestamp"
	headerKeyAuthorization = "Authorization"
	headerKeyContentType   = "Content-Type"
)

type client struct {
	baseURLPublic   string
	baseURLInternal string
	access          string
	secret          string
	username        string
	password        string
	httpClient      *http.Client

	jwt string
	mu  sync.Mutex
}

func NewClient(baseURL, access, secret, username, password string) *client {
	baseURLPublic := fmt.Sprintf("https://cloudapi.%s", baseURL)
	baseURLInternal := fmt.Sprintf("https://api.%s", baseURL)
	return &client{
		baseURLPublic:   baseURLPublic,
		baseURLInternal: baseURLInternal,
		access:          access,
		secret:          secret,
		httpClient:      http.DefaultClient,
		username:        username,
		password:        password,
	}
}

func (c *client) createRequest(ctx context.Context, method, uri string, body interface{}) (*http.Request, error) {
	var bb io.Reader

	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}

		bb = bytes.NewReader(b)
	}

	r, err := http.NewRequestWithContext(ctx, method, uri, bb)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r.Header.Add(headerKeyContentType, "application/json")

	return r, nil
}

func (c *client) Do(ctx context.Context, method, uri string, body interface{}) (*http.Response, error) {
	r, err := c.createRequest(ctx, method, c.baseURLPublic+uri, body)
	if err != nil {
		return nil, err
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)
	r.Header.Add(headerKeyTimestamp, now)

	payload := strings.Join([]string{method, uri, now}, "\n")
	h := hmac.New(sha256.New, []byte(c.secret))
	h.Write([]byte(payload))

	bearer := "Bearer " + c.access + ":" + base64.StdEncoding.EncodeToString(h.Sum(nil))
	r.Header.Add(headerKeyAuthorization, bearer)

	return c.httpClient.Do(r)
}

func (c *client) doInternal(ctx context.Context, method, uri string, body interface{}, forceRefreshJwt bool) (*http.Response, error) {
	r, err := c.createRequest(ctx, method, c.baseURLInternal+uri, body)
	if err != nil {
		return nil, err
	}

	jwt, err := c.getJwt(ctx, forceRefreshJwt)
	if err != nil {
		return nil, err
	}

	bearer := "Bearer " + jwt
	r.Header.Add(headerKeyAuthorization, bearer)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	// If we force refresh jwt and still get auth error then the credentials are probably wrong so prevent looping forever
	if resp.StatusCode == 401 && !forceRefreshJwt {
		return c.doInternal(ctx, method, uri, body, true)
	}

	return resp, nil
}

func (c *client) DoInternal(ctx context.Context, method, uri string, body interface{}) (*http.Response, error) {
	return c.doInternal(ctx, method, uri, body, false)
}

func (c *client) getJwt(ctx context.Context, forceRefresh bool) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.jwt == "" || forceRefresh {
		r, err := c.createRequest(ctx, "POST", c.baseURLInternal+sessionsPath, nil)
		if err != nil {
			return "", err
		}

		r.SetBasicAuth(c.username, c.password)

		res, err := c.httpClient.Do(r)
		if err != nil {
			return "", err
		}

		if res.StatusCode < 200 || res.StatusCode >= 300 {
			bb, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return "", fmt.Errorf("create session failed: reason could not be determined: %v", err)
			}
			return "", fmt.Errorf("create session failed: %s", string(bb))
		}

		bb, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		var respBody sessionsResponse
		if err := json.Unmarshal(bb, &respBody); err != nil {
			return "", err
		}

		c.jwt = respBody.Jwt
	}

	return c.jwt, nil
}
