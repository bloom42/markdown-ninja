package pingoo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type KeyType string
type JwkUse int32

type Algorithm string
type Type string
type Curve string

const (
	algorithmNone Algorithm = "none"

	AlgorithmHS256 Algorithm = "HS256"
	AlgorithmHS512 Algorithm = "HS512"
	AlgorithmEdDsa Algorithm = "EdDSA"
)

const (
	TypeJWT Type = "JWT"
)

const (
	CurveEd25519 Curve = "Ed25519"
)

const (
	KeyTypeOKP KeyType = "OKP"
)

const (
	JwkUseSign JwkUse = iota
	JwkUseEncrypt
)

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

func (client *Client) periodicallyFetchJkws() {
	for {
		err := client.fetchJkws()
		if err != nil {
			client.logger.Error("pingoo: error fetching JWKS", slog.String("error", err.Error()))
		}
		time.Sleep(client.jwksFetchInterval)
	}
}

func (client *Client) fetchJkws() error {
	var jwksRes Jwks

	err := client.request(context.Background(), requestParams{
		Method: http.MethodGet,
		Route:  fmt.Sprintf("/jwks/%s", client.projectId.String()),
	}, &jwksRes)
	if err != nil {
		return err
	}

	client.jwksLock.Lock()
	client.jwks = jwksRes
	// TODO: validate JWKS?
	client.jwksLock.Unlock()

	client.logger.Debug("pingoo: JWKS successfully refreshed")

	return nil
}
