package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func parseRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}

	// Validate PEM block type
	if block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM block type %q, expected \"RSA PRIVATE KEY\" or \"PRIVATE KEY\"", block.Type)
	}

	// Try PKCS#1 first
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		if key.N.BitLen() < 2048 {
			return nil, fmt.Errorf("RSA key too small (%d bits), minimum 2048", key.N.BitLen())
		}
		return key, nil
	}

	// Try PKCS#8
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key (tried PKCS#1 and PKCS#8): %w", err)
	}

	rsaKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an RSA private key")
	}
	if rsaKey.N.BitLen() < 2048 {
		return nil, fmt.Errorf("RSA key too small (%d bits), minimum 2048", rsaKey.N.BitLen())
	}
	return rsaKey, nil
}

func parseECPrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}

	// Validate PEM block type
	if block.Type != "EC PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM block type %q, expected \"EC PRIVATE KEY\" or \"PRIVATE KEY\"", block.Type)
	}

	// Try SEC 1 (EC) format first
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err == nil {
		if err := validateECCurve(key); err != nil {
			return nil, err
		}
		return key, nil
	}

	// Try PKCS#8
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	ecKey, ok := parsed.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an EC private key")
	}
	if err := validateECCurve(ecKey); err != nil {
		return nil, err
	}
	return ecKey, nil
}

func validateECCurve(key *ecdsa.PrivateKey) error {
	curve := key.Curve
	if curve == elliptic.P256() || curve == elliptic.P384() || curve == elliptic.P521() {
		return nil
	}
	return fmt.Errorf("unsupported EC curve %q, must be P-256, P-384, or P-521", curve.Params().Name)
}

func parseEdPrivateKey(data []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM block found")
	}

	// Validate PEM block type
	if block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected PEM block type %q, expected \"PRIVATE KEY\"", block.Type)
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ed25519 private key: %w", err)
	}

	edKey, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an Ed25519 private key")
	}
	if len(edKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid Ed25519 private key size: got %d bytes, expected %d", len(edKey), ed25519.PrivateKeySize)
	}
	return edKey, nil
}
