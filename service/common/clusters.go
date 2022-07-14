package common

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/couchbaselabs/cbdynclusterd/cluster"
	"github.com/couchbaselabs/cbdynclusterd/dyncontext"
	"github.com/couchbaselabs/cbdynclusterd/helper"
	"github.com/couchbaselabs/cbdynclusterd/service"
)

func NewNode(hostname string, connCtx service.ConnectContext) *Node {
	restUsername := helper.RestUser
	if connCtx.RestUsername != "" {
		restUsername = connCtx.RestUsername
	}
	restPassword := helper.RestPass
	if connCtx.RestPassword != "" {
		restPassword = connCtx.RestPassword
	}
	sshUsername := helper.SshUser
	if connCtx.SshUsername != "" {
		sshUsername = connCtx.SshUsername
	}
	sshPassword := helper.SshPass
	if connCtx.SshPassword != "" {
		sshPassword = connCtx.SshPassword
	}
	return &Node{
		HostName:  hostname,
		Port:      strconv.Itoa(helper.GetRestPort(connCtx.UseSecure)),
		SshLogin:  &helper.Cred{Username: sshUsername, Password: sshPassword, Hostname: hostname, Port: helper.SshPort, KeyPath: connCtx.SshKeyPath, Secure: connCtx.UseSecure},
		RestLogin: &helper.Cred{Username: restUsername, Password: restPassword, Hostname: hostname, Port: helper.GetRestPort(connCtx.UseSecure), Secure: connCtx.UseSecure},
		N1qlLogin: &helper.Cred{Username: restUsername, Password: restPassword, Hostname: hostname, Port: helper.GetN1qlPort(connCtx.UseSecure), Secure: connCtx.UseSecure},
		FtsLogin:  &helper.Cred{Username: restUsername, Password: restPassword, Hostname: hostname, Port: helper.GetFtsPort(connCtx.UseSecure), Secure: connCtx.UseSecure},
	}
}

func SetupCluster(opts ClusterSetupOptions, connCtx service.ConnectContext) (string, error) {
	services := opts.Services

	var nodes []*Node
	for i, service := range services {
		hostname := opts.Nodes[i].IPv4Address
		node := NewNode(hostname, connCtx)
		node.Services = service
		nodes = append(nodes, node)
	}

	config := Config{
		MemoryQuota:   opts.MemoryQuota,
		StorageMode:   opts.StorageMode,
		User:          opts.User,
		Bucket:        opts.Bucket,
		UseHostname:   opts.UseHostname,
		UseDevPreview: opts.UseDeveloperPreview,
	}

	clusterManager := &Manager{
		Nodes:  nodes,
		Config: config,
	}

	return clusterManager.StartCluster()
}

func ConnString(ctx context.Context, s service.ClusterService, clusterID string, useSSL bool) (string, error) {
	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return "", err
	}

	var addresses []string
	for _, node := range c.Nodes {
		addresses = append(addresses, node.IPv4Address)
	}

	scheme := "couchbase"
	if useSSL {
		scheme = "couchbases"
	}

	return fmt.Sprintf("%s://%s", scheme, strings.Join(addresses, ",")), nil
}

func AddSampleBucket(ctx context.Context, s service.ClusterService, clusterID string, opts service.AddSampleOptions, connCtx service.ConnectContext) error {
	log.Printf("Loading sample bucket %s to cluster %s (requested by: %s)", opts.SampleBucket, clusterID, dyncontext.ContextUser(ctx))

	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return err
	}

	if helper.SampleBucketsCount[opts.SampleBucket] == 0 {
		return errors.New("Unknown sample bucket")
	}

	if len(c.Nodes) == 0 {
		return errors.New("no nodes available")
	}

	node := NewNode(c.Nodes[0].IPv4Address, connCtx)

	return node.LoadSample(opts.SampleBucket)
}

func AddBucket(ctx context.Context, s service.ClusterService, clusterID string, opts service.AddBucketOptions, connCtx service.ConnectContext) error {
	log.Printf("Adding bucket %s to cluster %s (requested by: %s)", opts.Name, clusterID, dyncontext.ContextUser(ctx))

	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return err
	}

	if len(c.Nodes) == 0 {
		return errors.New("no nodes available")
	}

	node := NewNode(c.Nodes[0].IPv4Address, connCtx)

	return node.CreateBucket(&cluster.Bucket{
		Name:              opts.Name,
		Type:              opts.BucketType,
		ReplicaCount:      opts.ReplicaCount,
		RamQuotaMB:        strconv.Itoa(opts.RamQuota),
		EphEvictionPolicy: opts.EvictionPolicy,
		StorageBackend:    opts.StorageBackend,
	})
}

func SetupCertAuth(ctx context.Context, s service.ClusterService, clusterID string, opts service.SetupClientCertAuthOptions, connCtx service.ConnectContext) (*service.CertAuthResult, error) {
	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	initialNodes := c.Nodes
	clusterVersion := initialNodes[0].InitialServerVersion
	var nodes []Node
	for _, node := range initialNodes {
		nodes = append(nodes, *NewNode(node.IPv4Address, connCtx))
	}

	return setupCertAuth(opts.UserName, opts.UserEmail, nodes, clusterVersion, opts.NumRoots)
}

func SetupClusterEncryption(ctx context.Context, s service.ClusterService, clusterID string, opts service.SetupClusterEncryptionOptions, connCtx service.ConnectContext) error {
	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return err
	}

	var nodes []Node
	for _, node := range c.Nodes {
		nodes = append(nodes, *NewNode(node.IPv4Address, connCtx))
	}

	return setupClusterEncryption(nodes, opts)
}

func AddCollection(ctx context.Context, s service.ClusterService, clusterID string, opts service.AddCollectionOptions, connCtx service.ConnectContext) error {
	log.Printf("Adding collection %s to bucket %s on cluster %s (requested by: %s)", opts.Name,
		opts.BucketName, clusterID, dyncontext.ContextUser(ctx))

	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return err
	}

	if len(c.Nodes) == 0 {
		return errors.New("no nodes available")
	}

	node := NewNode(c.Nodes[0].IPv4Address, connCtx)

	return node.CreateCollection(&cluster.Collection{
		Name:       opts.Name,
		ScopeName:  opts.ScopeName,
		BucketName: opts.BucketName,
	})
}

func KillAllClusters(ctx context.Context, s service.ClusterService) error {
	clusters, err := s.GetAllClusters(ctx)
	if err != nil {
		return err
	}

	var clustersToKill []string

	for _, c := range clusters {
		clustersToKill = append(clustersToKill, c.ID)
	}

	signal := make(chan error)

	for _, clusterID := range clustersToKill {
		go func(clusterID string) {
			signal <- s.KillCluster(ctx, clusterID)
		}(clusterID)
	}

	var killError error
	for range clustersToKill {
		err := <-signal
		if err != nil && killError == nil {
			killError = err
		}
	}
	if killError != nil {
		return killError
	}

	return nil
}

func AddUser(ctx context.Context, s service.ClusterService, clusterID string, opts service.AddUserOptions, connCtx service.ConnectContext) error {
	log.Printf("Adding user %s on cluster %s (requested by: %s)", opts.User.Name, clusterID, dyncontext.ContextUser(ctx))

	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return err
	}

	if len(c.Nodes) == 0 {
		return errors.New("no nodes available")
	}

	node := NewNode(c.Nodes[0].IPv4Address, connCtx)

	return node.CreateUser(opts.User)
}

func RunCBCollect(ctx context.Context, s service.ClusterService, clusterID string, connCtx service.ConnectContext) (*service.CBCollectResult, error) {
	c, err := s.GetCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	type Resp struct {
		Err   error
		Bytes []byte
	}

	recv := make(chan Resp)

	var nodes []*Node
	for _, node := range c.Nodes {
		node := NewNode(node.IPv4Address, connCtx)
		nodes = append(nodes, node)
		go func(recv chan Resp, node *Node) {
			bytes, err := node.RunCBCollect()
			recv <- Resp{Err: err, Bytes: bytes}
		}(recv, node)
	}

	collections := make(map[string][]byte, len(nodes))

	for _, n := range nodes {
		resp := <-recv
		if resp.Err != nil {
			err = resp.Err
			continue
		}
		collections[n.HostName] = resp.Bytes
	}

	if err != nil {
		return nil, err
	}

	return &service.CBCollectResult{Collections: collections}, nil
}
