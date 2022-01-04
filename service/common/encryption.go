package common

import (
	"github.com/couchbaselabs/cbdynclusterd/service"
)

func setupClusterEncryption(nodes []Node, opts service.SetupClusterEncryptionOptions) error {
	epnode := nodes[0]

	autoFailoverEnabled, err := epnode.IsAutoFailoverEnabled()
	if err != nil {
		return err
	}

	// disable auto failover, you can't enable node to node encryption with it enabled
	if autoFailoverEnabled {
		if err = epnode.SetAutoFailover(false); err != nil {
			return err
		}
	}

	// setup external listeners, not too sure what this does but it is required
	for _, node := range nodes {
		if err = node.ManageExternalListener(true); err != nil {
			return err
		}
	}

	// change node encryption
	for _, node := range nodes {
		if err = node.EnableClusterEncryption(); err != nil {
			return err
		}
	}

	// disable external listener
	for _, node := range nodes {
		if err = node.ManageExternalListener(false); err != nil {
			return err
		}
	}

	// change level
	if opts.Level != "" {
		version, err := getVersion(&epnode)
		if err != nil {
			return err
		}
		// Strict encryption level is hidden behind a flag in 6.6.5
		if opts.Level == "strict" && version.Major == 6 {
			if err := epnode.AllowStrictEncryption(); err != nil {
				return err
			}
		}
		if err := epnode.SetClusterEncryptionLevel(opts.Level); err != nil {
			return err
		}
	}

	// enable auto failover again if it was originally enabled
	if autoFailoverEnabled {
		if err = epnode.SetAutoFailover(true); err != nil {
			return err
		}
	}

	return nil
}
