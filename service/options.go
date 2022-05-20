package service

import "github.com/couchbaselabs/cbdynclusterd/helper"

type AddCollectionOptions struct {
	Name       string
	ScopeName  string
	BucketName string
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
	ReplicaCount   int
	BucketType     string
	EvictionPolicy string
	StorageBackend string
}

type AddSampleOptions struct {
	SampleBucket string
}

type SetupClusterEncryptionOptions struct {
	Level string
}

type AddUserOptions struct {
	User *helper.UserOption
}
