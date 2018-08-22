package daemon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type NodeOptions struct {
	Name          string
	Platform      string
	ServerVersion string
}

func allocateNode(ctx context.Context, clusterID string, timeout time.Time, opts NodeOptions) (string, error) {
	log.Printf("Allocating node for cluster %s (requested by: %s)", clusterID, ContextUser(ctx))

	containerName := fmt.Sprintf("dynclsr-%s-%s", clusterID, opts.Name)

	// TODO(brett19): This needs to be able to understand server versions better
	containerImage := fmt.Sprintf("couchbase/server:%s", opts.ServerVersion)

	createResult, err := docker.ContainerCreate(context.Background(), &container.Config{
		Image: containerImage,
		Labels: map[string]string{
			"com.couchbase.dyncluster.creator":                ContextUser(ctx),
			"com.couchbase.dyncluster.cluster_id":             clusterID,
			"com.couchbase.dyncluster.node_name":              opts.Name,
			"com.couchbase.dyncluster.initial_server_version": opts.ServerVersion,
		},
	}, &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: "macvlan0",
	}, nil, containerName)
	if err != nil {
		return "", err
	}

	err = docker.ContainerStart(context.Background(), createResult.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return createResult.ID, nil
}

func killNode(ctx context.Context, containerID string) error {
	log.Printf("Killing node %s (requested by: %s)", containerID, ContextUser(ctx))

	err := docker.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return err
	}

	// No need to kill the node, since we use `kill on stop` when creating the account
	/*
		err = docker.ContainerKill(context.Background(), containerID, "")
		if err != nil {
			return err
		}
	*/

	return nil
}
