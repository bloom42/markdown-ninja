package jwt

// JSON Web Keyset
type Jwks struct {
	Keys []Jwk `json:"keys"`
}

// JSON Web Key
// TODO: validate when unmarshalling from JSON
type Jwk struct {
	KeyID     string    `json:"kid"`
	KeyType   KeyType   `json:"kty"`
	Algorithm Algorithm `json:"alg"`
	// #[serde(flatten)]
	// pub crypto: JwtKeyCrypto,
	Curve Curve             `json:"crv"`
	X     BytesBase64RawUrl `json:"x"`
	Use   string            `json:"use"`
}
