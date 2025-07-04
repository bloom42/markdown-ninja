package pingoo

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

var (
	ErrTokenIsNotValid     = errors.New("The token is not valid")
	ErrSignatureIsNotValid = errors.New("Signature is not valid")
	ErrTokenHasExpired     = errors.New("The token has expired")
	ErrAlgorithmIsNotValid = fmt.Errorf("Algorithm is not valid. Valid algorithms values are: [%s, %s]", AlgorithmHS256, AlgorithmHS512)
	ErrKeyNotFound         = func(keyID string) error {
		return fmt.Errorf("JWT key not found: %s", keyID)
	}
	ErrIssuerIsNotValid = errors.New("issuer is not valud")
)

type header struct {
	Algorithm Algorithm `json:"alg"`
	Type      Type      `json:"typ"`
	KeyID     string    `json:"kid"`
}

type VerifyingKey interface {
	Verify(message, signature []byte) bool
	Algorithm() Algorithm
	ID() string
}

// type JwksClientConfig struct {
// 	Endpoint      string
// 	FetchInterval time.Duration
// 	HttpClient    *http.Client
// 	Logger        *slog.Logger
// 	UserAgent     string
// }

// type JwksClient struct {
// 	endpoint      string
// 	fetchInterval time.Duration
// 	httpClient    *http.Client
// 	logger        *slog.Logger
// 	userAgent     string

// 	// a map of verifying keys, by KeyID
// 	keys     map[string]VerifyingKey
// 	jwksKeysLock sync.RWMutex
// }

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

type Ed25519VerifyingKey struct {
	keyID     string
	publicKey ed25519.PublicKey
}

func NewEd25519VerifyingKey(keyID string, publicKey []byte) (Ed25519VerifyingKey, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return Ed25519VerifyingKey{}, errors.New("invalid ed25519 public key size")
	}

	return Ed25519VerifyingKey{keyID: keyID, publicKey: ed25519.PublicKey(publicKey)}, nil
}

func (key Ed25519VerifyingKey) Algorithm() Algorithm {
	return AlgorithmEdDsa
}

func (key Ed25519VerifyingKey) ID() string {
	return key.keyID
}

func (key Ed25519VerifyingKey) Verify(message, signature []byte) bool {
	return ed25519.Verify(key.publicKey, message, signature)
}

// func NewJwksClient(config JwksClientConfig) *JwksClient {
// 	fetchInterval := config.FetchInterval
// 	if fetchInterval == 0 {
// 		fetchInterval = time.Minute
// 	}

// 	httpClient := config.HttpClient
// 	if httpClient == nil {
// 		httpClient = httpx.DefaultClient()
// 	}

// 	client := &JwksClient{
// 		endpoint:      config.Endpoint,
// 		fetchInterval: fetchInterval,
// 		httpClient:    httpClient,
// 		logger:        config.Logger,
// 		userAgent:     config.UserAgent,

// 		keys:     make(map[string]VerifyingKey),
// 		jwksKeysLock: sync.RWMutex{},
// 	}
// 	go client.periodicallyFetchJkws()
// 	return client
// }

// Usage: client.Verify(jwt, &myClaimsStruct)
func (client *Client) VerifyJWT(token string, claimsStructPointer any) (err error) {
	if strings.Count(token, ".") != 2 {
		return ErrTokenIsNotValid
	}

	// Header
	var header header
	headerEnd := strings.IndexByte(token, '.')
	encodedHeader := token[:headerEnd]
	headerJson, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return ErrTokenIsNotValid
	}
	err = json.Unmarshal(headerJson, &header)
	if err != nil {
		return ErrTokenIsNotValid
	}

	client.jwksKeysLock.RLock()
	key, keyOk := client.jwksKeys[header.KeyID]
	if !keyOk {
		// if key not found, we first try to fetch the JWKS again
		client.fetchJkws()
		key, keyOk = client.jwksKeys[header.KeyID]
		if !keyOk {
			client.jwksKeysLock.RUnlock()
			return ErrKeyNotFound(header.KeyID)
		}
	}
	client.jwksKeysLock.RUnlock()

	if header.Algorithm != key.Algorithm() || header.Type != TypeJWT {
		return ErrTokenIsNotValid
	}

	// Signature
	signatureStart := strings.LastIndexByte(token, '.')
	encodedSignature := token[signatureStart+1:]
	signature, err := base64.RawURLEncoding.DecodeString(encodedSignature)
	if err != nil {
		return ErrTokenIsNotValid
	}

	encodedHeaderAndClaims := token[:signatureStart]
	if !key.Verify([]byte(encodedHeaderAndClaims), signature) {
		return ErrTokenIsNotValid
	}

	// Claims
	encodedClaims := token[headerEnd+1 : signatureStart]
	claimsJson, err := base64.RawURLEncoding.DecodeString(encodedClaims)
	if err != nil {
		return ErrTokenIsNotValid
	}

	err = json.Unmarshal(claimsJson, claimsStructPointer)
	if err != nil {
		return ErrTokenIsNotValid
	}

	return nil
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

	ed25519Keys := make([]Ed25519VerifyingKey, 0, len(jwksRes.Keys))
	for _, key := range jwksRes.Keys {
		if key.Algorithm != AlgorithmEdDsa {
			client.logger.Error(fmt.Sprintf("pingoo: unsuported JWK algorithm (keyID: %s): %s", key.KeyID, key.Algorithm))
			continue
		}
		verifyingKey, err := NewEd25519VerifyingKey(key.KeyID, key.X)
		if err != nil {
			return err
		}
		ed25519Keys = append(ed25519Keys, verifyingKey)
	}

	client.jwksKeysLock.Lock()
	client.jwksKeys = make(map[string]VerifyingKey, len(ed25519Keys))
	for _, key := range ed25519Keys {
		client.jwksKeys[key.ID()] = key
	}
	client.jwksKeysLock.Unlock()

	client.logger.Debug("pingoo: JWKS successfully refreshed")

	return nil
}
