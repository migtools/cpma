package ocp

type Containers struct {
    Registries map[string]Registries
}

type Registries struct {
    List []string `toml:"registries"`
}

type ImageCR struct {
    APIVersion string   `yaml:"apiVersion"`
    Kind       string   `yaml:"kind"`
    Metadata   Metadata `yaml:"metadata"`
    Spec       struct {
        RegistrySources RegistrySources `yaml:"registrySources"`
    } `yaml:"spec"`
}

type Metadata struct {
    Name        string
    Annotations map[string]string `yaml:"annotations"`
}

type RegistrySources struct {
    BlockedRegistries  []string `yaml:"blockedRegistries,omitempty"`
    InsecureRegistries []string `yaml:"insecureRegistries,omitempty"`
}

type RegistriesTransform struct {
    Config *Config
}
