package ocp

type NetworkCR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		ClusterNetworks []ClusterNetwork `yaml:"clusterNetwork"`
		ServiceNetwork  string           `yaml:"serviceNetwork"`
		DefaultNetwork  `yaml:"defaultNetwork"`
	} `yaml:"spec"`
}

// ClusterNetwork contains CIDR and address size to assign to each node
type ClusterNetwork struct {
	CIDR       string `yaml:"cidr"`
	HostPrefix uint32 `yaml:"hostPrefix"`
}

// DefaultNetwork containts network type and SDN plugin name
type DefaultNetwork struct {
	Type               string `yaml:"type"`
	OpenshiftSDNConfig struct {
		Mode string `yaml:"mode"`
	} `yaml:"openshiftSDNConfig"`
}

type SDNTransform struct {
	Config *Config
}
