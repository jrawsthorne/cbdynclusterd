package docker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/couchbaselabs/cbdynclusterd/dyncontext"
	"github.com/couchbaselabs/cbdynclusterd/service/common"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

var NetworkName = "macvlan0"

func (ds *DockerService) allocateNode(ctx context.Context, clusterID string, timeout time.Time, opts common.NodeOptions) (string, error) {
	log.Printf("Allocating node for cluster %s (requested by: %s)", clusterID, dyncontext.ContextUser(ctx))

	containerName := fmt.Sprintf("dynclsr-%s-%s", clusterID, opts.Name)
	containerImage := opts.VersionInfo.ToImageName(ds.dockerRegistry)

	createResult, err := ds.docker.ContainerCreate(context.Background(), &container.Config{
		Image: containerImage,
		Labels: map[string]string{
			"com.couchbase.dyncluster.creator":                dyncontext.ContextUser(ctx),
			"com.couchbase.dyncluster.cluster_id":             clusterID,
			"com.couchbase.dyncluster.node_name":              opts.Name,
			"com.couchbase.dyncluster.initial_server_version": opts.ServerVersion,
		},
		// same effect as ntp
		Volumes: map[string]struct{}{"/etc/localtime:/etc/localtime": {}},
	}, &container.HostConfig{
		AutoRemove:  true,
		NetworkMode: container.NetworkMode(NetworkName),
		CapAdd:      []string{"NET_ADMIN"},
	}, nil, containerName)
	if err != nil {
		return "", err
	}

	err = ds.docker.ContainerStart(context.Background(), createResult.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return createResult.ID, nil
}

func (ds *DockerService) killNode(ctx context.Context, containerID string) error {
	log.Printf("Killing node %s (requested by: %s)", containerID, dyncontext.ContextUser(ctx))

	err := ds.docker.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return err
	}

	// No need to kill the node, since we use `kill on stop` when creating the container
	/*
		err = docker.ContainerKill(context.Background(), containerID, "")
		if err != nil {
			return err
		}
	*/

	return nil
}
