package configmaps

// ConfigMap represent configmap definition
type ConfigMap struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   MetaData `json:"metadata"`
	Data       Data     `json:"data"`
}

// MetaData configmap's metadata
type MetaData struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// Data contains CA
type Data struct {
	CAData string `json:"ca"`
}

const (
	// APIVersion is the apiVersion string
	APIVersion = "v1"
	// Kind is config map resource type
	Kind = "ConfigMap"
)

// GenConfigMap generates a secret
func GenConfigMap(name string, namespace string, CAData []byte) *ConfigMap {
	return &ConfigMap{
		APIVersion: APIVersion,
		Data: Data{
			CAData: string(CAData),
		},
		Kind: Kind,
		Metadata: MetaData{
			Name:      name,
			Namespace: namespace,
		},
	}
}
