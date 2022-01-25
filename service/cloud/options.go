package cloud

import (
	"errors"
)

type NodeSetupOptions struct {
	Services []V3CouchbaseService
	Size     uint32
}

type ClusterSetupOptions struct {
	Nodes []NodeSetupOptions
}

func ParseServiceName(service string) (V3CouchbaseService, error) {
	switch service {
	case "kv":
		return V3CouchbaseServiceData, nil
	case "index":
		return V3CouchbaseServiceIndex, nil
	case "cbas":
		return V3CouchbaseServiceAnalytics, nil
	case "n1ql":
		return V3CouchbaseServiceQuery, nil
	case "fts":
		return V3CouchbaseServiceSearch, nil
	case "eventing":
		return V3CouchbaseServiceEventing, nil
	default:
		return V3CouchbaseServiceData, errors.New("Unknown service")
	}
}
