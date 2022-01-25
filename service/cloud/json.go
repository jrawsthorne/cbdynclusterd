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
	V3StorageTypeGP3 V3StorageType = "GP3"
	V3StorageTypeIO2 V3StorageType = "IO2"
)

var defaultCompute = "m5.xlarge"
var defaultRegion = "us-west-2"

type V3ServersStorage struct {
	Type V3StorageType `json:"type"`
	IOPS uint32        `json:"IOPS"`
	Size uint32        `json:"size"`
}

var defaultStorage = V3ServersStorage{
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

type databaseUserJSON struct {
	Username         string       `json:"username"`
	Password         string       `json:"password"`
	AllBucketsAccess V3BucketRole `json:"allBucketsAccess"`
}

type bucketSpecJSON struct {
	Name        string `json:"name"`
	MemoryQuota int    `json:"memoryQuota"`
	Replicas    int    `json:"replicas"`
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
