package cloud

type getAllClustersClusterJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getAllClustersJSON struct {
	Data struct {
		Items []getAllClustersClusterJSON `json:"items"`
	} `json:"data"`
}

type getClusterJSONVersion struct {
	Name string `json:"name"`
}

type getClusterJSON struct {
	ID           string                `json:"id"`
	ProjectID    string                `json:"projectId"`
	Status       string                `json:"status"`
	Name         string                `json:"name"`
	Version      getClusterJSONVersion `json:"version"`
	EndpointsSRV string                `json:"endpointsSrv"`
}

type awsServerJSON struct {
	InstanceSize string `json:"instanceSize"`
	EbsSizeGib   int    `json:"ebsSizeGib"`
}

type nodeJson struct {
	Services []string      `json:"services"`
	Size     uint32        `json:"size"`
	AWS      awsServerJSON `json:"aws"`
}

type V3StorageType string

const (
	V3StorageTypeGP3    V3StorageType = "GP3"
	V3StorageTypeIO2    V3StorageType = "IO2"
	V3StorageTypePD_SSD V3StorageType = "PD-SSD"
)

type V3ServersStorage struct {
	Type V3StorageType `json:"type"`
	IOPS uint32        `json:"IOPS,omitempty"`
	Size uint32        `json:"size"`
}

var defaultComputeAWS = "m5.xlarge"
var defaultRegionAWS = "us-east-2"
var defaultProvider = V3ProviderAWS
var defaultSingleAZ = true
var defaultComputeGCP = "n2-standard-4"
var defaultRegionGCP = "us-east1"
var defaultStorageGCP = V3ServersStorage{
	Type: V3StorageTypePD_SSD,
	Size: 50,
}
var defaultStorageAWS = V3ServersStorage{
	Type: V3StorageTypeGP3,
	IOPS: 3000,
	Size: 50,
}

type V3CouchbaseService string

const (
	V3CouchbaseServiceData      V3CouchbaseService = "data"
	V3CouchbaseServiceIndex     V3CouchbaseService = "index"
	V3CouchbaseServiceQuery     V3CouchbaseService = "query"
	V3CouchbaseServiceSearch    V3CouchbaseService = "search"
	V3CouchbaseServiceEventing  V3CouchbaseService = "eventing"
	V3CouchbaseServiceAnalytics V3CouchbaseService = "analytics"
)

type V3Server struct {
	Size     uint32               `json:"size"`
	Compute  string               `json:"compute"`
	Services []V3CouchbaseService `json:"services"`
	Storage  V3ServersStorage     `json:"storage"`
}

type supportPackageJson struct {
	Type     string `json:"type"`
	Timezone string `json:"timezone"`
}

var defaultSupportPackage = supportPackageJson{
	Type:     "DeveloperPro",
	Timezone: "PT",
}

var defaultAWS = awsServerJSON{
	InstanceSize: "m5.xlarge",
	EbsSizeGib:   50,
}

type V3Provider string

const (
	V3ProviderAWS   V3Provider = "aws"
	V3ProviderAzure V3Provider = "azure"
	V3ProviderGCP   V3Provider = "gcp"
)

type V3PlaceHosted struct {
	Provider V3Provider `json:"provider"`
	Region   string     `json:"region"`
	CIDR     string     `json:"CIDR"`
}

type V3Place struct {
	SingleAZ bool          `json:"singleAZ"`
	Hosted   V3PlaceHosted `json:"hosted"`
}

type V3Environment string

const (
	V3EnvironmentHosted V3Environment = "hosted"
	V3EnvironmentVPC    V3Environment = "vpc"
)

type setupClusterJson struct {
	Environment    V3Environment      `json:"environment"`
	Name           string             `json:"clusterName"`
	ProjectID      string             `json:"projectId"`
	Place          V3Place            `json:"place"`
	Servers        []V3Server         `json:"servers"`
	SupportPackage supportPackageJson `json:"supportPackage"`
}

type V3BucketRole string

const (
	V3BucketRoleDataReader V3BucketRole = "data_reader"
	V3BucketRoleDataWriter V3BucketRole = "data_writer"
)

type ScopePermission struct {
	Name string `json:"name"`
}

type BucketPermission struct {
	Name   string            `json:"name"`
	Scopes []ScopePermission `json:"scopes,omitempty"`
}

type Permissions struct {
	Buckets []BucketPermission `json:"buckets,omitempty"`
}

type ReadWritePermissions struct {
	DataReader *Permissions `json:"data_reader"`
	DataWriter *Permissions `json:"data_writer"`
}

type databaseUserJSON struct {
	Name        string               `json:"name"`
	Password    string               `json:"password"`
	Permissions ReadWritePermissions `json:"permissions"`
}

// TODO: Use public API when AV-27634 fixed

// type databaseUserJSON struct {
// 	Username         string       `json:"username"`
// 	Password         string       `json:"password"`
// 	AllBucketsAccess V3BucketRole `json:"allBucketsAccess"`
// }

type bucketSpecJSON struct {
	Name              string `json:"name"`
	MemoryQuota       int    `json:"memoryAllocationInMb"`
	Replicas          int    `json:"replicas"`
	ConflicResolution string `json:"bucketConflictResolution"`
	DurabilityLevel   string `json:"durabilityLevel"`
}

type bucketDeleteJSON struct {
	Name string `json:"name"`
}

type clusterHealthResponse struct {
	BucketStats struct {
		HealthStats    map[string]string `json:"healthStats"`
		Status         string            `json:"status"`
		TotalCount     int               `json:"totalCount"`
		UnhealthyCount int               `json:"unhealthyCount"`
	} `json:"bucketStats"`
}

type sessionsResponse struct {
	Jwt string `json:"jwt"`
}

type allowListJSON struct {
	CIDR    string `json:"cidr"`
	Comment string `json:"comment"`
}

type allowListBulkJSON struct {
	Create []allowListJSON `json:"create"`
}

type addSampleBucketJSON struct {
	Name string `json:"name"`
}

type getNodesJson struct {
	Data []struct {
		Data struct {
			Hostname string `json:"hostname"`
			ID       string `json:"id"`
		} `json:"data"`
	} `json:"data"`
}
