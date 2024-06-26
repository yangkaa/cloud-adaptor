package rke2

import (
	"goodrain.com/cloud-adaptor/internal/model"
	"testing"
)

func Test(t *testing.T) {
	installRKE2Cluster(model.RKE2Nodes{
		Role:       "agent",
		Host:       "8.130.119.244",
		Port:       22,
		User:       "root",
		Pass:       "gr123465!",
		ConfigFile: ``,
	})
}
