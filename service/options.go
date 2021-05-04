package service

type AddCollectionOptions struct {
	Name        string
	ScopeName   string
	BucketName  string
	UseHostname bool
}

type SetupClientCertAuthOptions struct {
	UserName  string
	UserEmail string
	NumRoots  int
}

type AddBucketOptions struct {
	Name           string
	StorageMode    string
	RamQuota       int
	UseHostname    bool
	ReplicaCount   int
	BucketType     string
	EvictionPolicy string
	StorageBackend string
}

type AddSampleOptions struct {
	SampleBucket string
	UseHostname  bool
}
