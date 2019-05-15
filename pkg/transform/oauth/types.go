package oauth

// CA secret name
type CA struct {
	Name string `yaml:"name"`
}

// TLSClientCert secret name
type TLSClientCert struct {
	Name string `yaml:"name"`
}

// TLSClientKey secret name
type TLSClientKey struct {
	Name string `yaml:"name"`
}

// ClientSecret is a client secret for a privuder
type ClientSecret struct {
	Name string `yaml:"name"`
}
