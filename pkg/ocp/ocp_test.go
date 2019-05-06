package ocp

import (
	"testing"

	"github.com/fusor/cpma/env"
	"github.com/stretchr/testify/assert"
)

func TestAddConfigMaster(t *testing.T) {
	// Init config with default master config paths
	ocpMaster := ConfigMaster{}
	ocpMaster.Add("example.com")

	assert.Equal(t, ConfigMaster{
		ConfigFile: ConfigFile{
			Hostname: "example.com",
			Path:     "/etc/origin/master/master-config.yaml",
		},
	}, ocpMaster)

	// Init config with different master config path
	env.Config().Set("MasterConfigFile", "/test/path/master.yml")
	ocpMaster = ConfigMaster{}
	ocpMaster.Add("example.com")

	assert.Equal(t, ConfigMaster{
		ConfigFile: ConfigFile{
		Hostname: "example.com",
		Path:     "/test/path/master.yml",
		},
	}, ocpMaster)
}

func TestAddConfigNode(t *testing.T) {
	// Init config with default node config paths
	ocpNode := ConfigNode{}
	ocpNode.Add("example.com")

	assert.Equal(t, ConfigNode{
		ConfigFile: ConfigFile{
		Hostname: "example.com",
		Path:     "/etc/origin/node/node-config.yaml",
		},
	}, ocpNode)

	// Init config with different node config paths
	env.Config().Set("NodeConfigFile", "/test/path/node.yml")
	ocpNode = ConfigNode{}
	ocpNode.Add("example.com")

	assert.Equal(t, ConfigNode{
		ConfigFile: ConfigFile{
		Hostname: "example.com",
		Path:     "/test/path/node.yml",
		},
	}, ocpNode)
}
