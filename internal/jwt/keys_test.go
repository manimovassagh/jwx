package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

// helper: generate RSA PKCS#1 PEM
func generateRSAPKCS1PEM(t *testing.T) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der := x509.MarshalPKCS1PrivateKey(key)
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})
}

// helper: generate RSA PKCS#8 PEM
func generateRSAPKCS8PEM(t *testing.T) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal PKCS8: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

// helper: generate EC SEC1 PEM
func generateECSEC1PEM(t *testing.T, curve elliptic.Curve) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal EC key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
}

// helper: generate Ed25519 PKCS#8 PEM
func generateEdPKCS8PEM(t *testing.T) []byte {
	t.Helper()
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate Ed25519 key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		t.Fatalf("failed to marshal Ed25519 key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

func TestParseRSAPrivateKey_PKCS1(t *testing.T) {
	pemData := generateRSAPKCS1PEM(t)
	key, err := parseRSAPrivateKey(pemData)
	if err != nil {
		t.Fatalf("parseRSAPrivateKey(PKCS1) unexpected error: %v", err)
	}
	if key == nil {
		t.Fatal("expected non-nil key")
	}
	if key.N.BitLen() < 2048 {
		t.Errorf("expected at least 2048-bit key, got %d", key.N.BitLen())
	}
}

func TestParseRSAPrivateKey_PKCS8(t *testing.T) {
	pemData := generateRSAPKCS8PEM(t)
	key, err := parseRSAPrivateKey(pemData)
	if err != nil {
		t.Fatalf("parseRSAPrivateKey(PKCS8) unexpected error: %v", err)
	}
	if key == nil {
		t.Fatal("expected non-nil key")
	}
}

func TestParseRSAPrivateKey_NoPEMBlock(t *testing.T) {
	_, err := parseRSAPrivateKey([]byte("this is not a PEM block"))
	if err == nil {
		t.Fatal("expected error for non-PEM input")
	}
	if got := err.Error(); got != "no PEM block found" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseRSAPrivateKey_WrongKeyType(t *testing.T) {
	// Encode an EC key as PKCS#8 and try to parse as RSA
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key as PKCS8: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	_, err = parseRSAPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for EC key parsed as RSA")
	}
	if got := err.Error(); got != "key is not an RSA private key" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseECPrivateKey_SEC1(t *testing.T) {
	pemData := generateECSEC1PEM(t, elliptic.P256())
	key, err := parseECPrivateKey(pemData)
	if err != nil {
		t.Fatalf("parseECPrivateKey(SEC1) unexpected error: %v", err)
	}
	if key == nil {
		t.Fatal("expected non-nil key")
	}
	if key.Curve != elliptic.P256() {
		t.Errorf("expected P-256 curve, got %s", key.Curve.Params().Name)
	}
}

func TestParseECPrivateKey_InvalidPEM(t *testing.T) {
	_, err := parseECPrivateKey([]byte("not a pem"))
	if err == nil {
		t.Fatal("expected error for non-PEM input")
	}
	if got := err.Error(); got != "no PEM block found" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseECPrivateKey_WrongKeyType(t *testing.T) {
	// Encode an RSA key as PKCS#8 and try to parse as EC
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		t.Fatalf("failed to marshal RSA key as PKCS8: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	_, err = parseECPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for RSA key parsed as EC")
	}
	if got := err.Error(); got != "key is not an EC private key" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseEdPrivateKey_Valid(t *testing.T) {
	pemData := generateEdPKCS8PEM(t)
	key, err := parseEdPrivateKey(pemData)
	if err != nil {
		t.Fatalf("parseEdPrivateKey unexpected error: %v", err)
	}
	if len(key) != ed25519.PrivateKeySize {
		t.Errorf("expected key length %d, got %d", ed25519.PrivateKeySize, len(key))
	}
}

func TestParseEdPrivateKey_InvalidPEM(t *testing.T) {
	_, err := parseEdPrivateKey([]byte("garbage"))
	if err == nil {
		t.Fatal("expected error for non-PEM input")
	}
	if got := err.Error(); got != "no PEM block found" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseEdPrivateKey_WrongKeyType(t *testing.T) {
	// Encode an RSA key as PKCS#8 and try to parse as Ed25519
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(rsaKey)
	if err != nil {
		t.Fatalf("failed to marshal RSA key as PKCS8: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	_, err = parseEdPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for RSA key parsed as Ed25519")
	}
	if got := err.Error(); got != "key is not an Ed25519 private key" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestParseECPrivateKey_PKCS8(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key as PKCS8: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	key, err := parseECPrivateKey(pemData)
	if err != nil {
		t.Fatalf("parseECPrivateKey(PKCS8) unexpected error: %v", err)
	}
	if key.Curve != elliptic.P384() {
		t.Errorf("expected P-384 curve, got %s", key.Curve.Params().Name)
	}
}

func TestParseRSAPrivateKey_InvalidDER(t *testing.T) {
	// Valid PEM block but garbage DER content
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("not-valid-der")})
	_, err := parseRSAPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for invalid DER content")
	}
}

func TestParseECPrivateKey_InvalidDER(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("not-valid-der")})
	_, err := parseECPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for invalid DER content")
	}
}

func TestParseEdPrivateKey_InvalidDER(t *testing.T) {
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: []byte("not-valid-der")})
	_, err := parseEdPrivateKey(pemData)
	if err == nil {
		t.Fatal("expected error for invalid DER content")
	}
}
