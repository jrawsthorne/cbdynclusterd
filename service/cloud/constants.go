package cloud

const (
	// public API
	createClusterPath  = "/v3/clusters"
	deleteClusterPath  = "/v3/clusters/"
	getClusterPath     = "/v3/clusters/"
	getAllClustersPath = "/v3/clusters"
	// TODO: Use public API when AV-27634 fixed
	// createUserPath     = "/v3/clusters/%s/users"
	clustersHealthPath = "/v3/clusters/%s/health"

	// internal API
	internalBasePath    = "/v2/organizations/%s/projects/%s/clusters/%s"
	createBucketPath    = internalBasePath + "/buckets"
	addIPPath           = internalBasePath + "/allowlists-bulk"
	addSampleBucketPath = internalBasePath + "/buckets/samples"
	getNodesPath        = internalBasePath + "/nodes"
	sessionsPath        = "/sessions"
	createUserPath      = internalBasePath + "/users"

	clusterHealthy          = "healthy"
	clusterDeleting         = "destroying"
	clusterDeploying        = "deploying"
	clusterDegraded         = "degraded"
	clusterDraft            = "draft"
	clusterDeploymentFailed = "deploymentFailed"
)
