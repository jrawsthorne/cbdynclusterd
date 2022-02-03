package service

type AddCollectionOptions struct {
	Name        string
	ScopeName   string
	BucketName  string
	UseHostname bool
	UseSecure   bool
}

type SetupClientCertAuthOptions struct {
	UserName    string
	UserEmail   string
	NumRoots    int
	SSHKeyPath  string
	SSHUsername string
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
	UseSecure      bool
}

type AddSampleOptions struct {
	SampleBucket string
	UseHostname  bool
	UseSecure    bool
}

type SetupClusterEncryptionOptions struct {
	Level     string
	UseSecure bool
}
