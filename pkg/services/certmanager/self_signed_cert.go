package certmanager

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

func generateSelfSignedCert() (*tls.Certificate, error) {
	var privateKey any
	var err error

	privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// _, privateKey, err = ed25519.GenerateKey(rand.Reader)
	if err != nil {
		err = fmt.Errorf("certmanager: generating private key: %w", err)
		return nil, err
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Not before the beginning of the current day
	now := time.Now().UTC()
	notBefore := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	notAfter := notBefore.Add(365 * 24 * time.Hour)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("certmanager: failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		NotBefore:    notBefore,
		NotAfter:     notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	if err != nil {
		return nil, fmt.Errorf("certmanager: Failed to create certificate: %w", err)
	}

	x509Cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("certmanager: Unable to marshal private key: %w", err)
	}
	x509Key := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})

	tlsCert, err := tls.X509KeyPair(x509Cert, x509Key)
	if err != nil {
		return nil, fmt.Errorf("certmanager: error loading x509 TLS certificate: %w", err)
	}

	return &tlsCert, nil
}

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}
