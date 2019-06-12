package oauth

// CA secret name
type CA struct {
	Name string `json:"name"`
}

// TLSClientCert secret name
type TLSClientCert struct {
	Name string `json:"name"`
}

// TLSClientKey secret name
type TLSClientKey struct {
	Name string `json:"name"`
}

// ClientSecret is a client secret for a privuder
type ClientSecret struct {
	Name string `json:"name"`
}
