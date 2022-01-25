package cloud

const (
	createClusterPath  = "/v3/clusters"
	deleteClusterPath  = "/v3/clusters/"
	getClusterPath     = "/v3/clusters/"
	getAllClustersPath = "/v3/clusters"
	// createBucketPath   = "/v2/clusters/%s/buckets"
	createUserPath = "/v3/clusters/%s/users"
	// addIPPath          = "/v2/clusters/%s/allowlist"
	// deleteBucketPath   = "/v2/clusters/%s/buckets"
	clustersHealthPath = "/v3/clusters/%s/health"

	clusterHealthy          = "healthy"
	clusterDeleting         = "destroying"
	clusterDeploying        = "deploying"
	clusterDegraded         = "degraded"
	clusterDraft            = "draft"
	clusterDeploymentFailed = "deploymentFailed"
)
