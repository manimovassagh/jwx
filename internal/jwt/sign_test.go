package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSign_HMAC(t *testing.T) {
	tests := []struct {
		name    string
		opts    SignOptions
		wantErr bool
	}{
		{
			name: "HS256 basic",
			opts: SignOptions{
				Algorithm: "HS256",
				Secret:    "mysecret",
				Claims:    `{"sub":"1234","name":"test"}`,
			},
			wantErr: false,
		},
		{
			name: "HS384",
			opts: SignOptions{
				Algorithm: "HS384",
				Secret:    "mysecret",
				Claims:    `{"sub":"1234"}`,
			},
			wantErr: false,
		},
		{
			name: "HS512",
			opts: SignOptions{
				Algorithm: "HS512",
				Secret:    "mysecret",
				Claims:    `{"sub":"1234"}`,
			},
			wantErr: false,
		},
		{
			name: "missing secret",
			opts: SignOptions{
				Algorithm: "HS256",
				Secret:    "",
				Claims:    `{"sub":"1234"}`,
			},
			wantErr: true,
		},
		{
			name: "invalid claims JSON",
			opts: SignOptions{
				Algorithm: "HS256",
				Secret:    "mysecret",
				Claims:    `not json`,
			},
			wantErr: true,
		},
		{
			name: "unsupported algorithm",
			opts: SignOptions{
				Algorithm: "FAKE256",
				Secret:    "mysecret",
				Claims:    `{"sub":"1234"}`,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := Sign(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				parts := strings.Split(token, ".")
				if len(parts) != 3 {
					t.Errorf("expected 3 parts in token, got %d", len(parts))
				}
			}
		})
	}
}

func TestSign_Roundtrip(t *testing.T) {
	token, err := Sign(SignOptions{
		Algorithm: "HS256",
		Secret:    "test-secret",
		Claims:    `{"sub":"roundtrip","name":"Test User","iat":1516239022}`,
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	decoded, err := Decode(token)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if decoded.Payload["sub"] != "roundtrip" {
		t.Errorf("expected sub=roundtrip, got %v", decoded.Payload["sub"])
	}
	if decoded.Payload["name"] != "Test User" {
		t.Errorf("expected name=Test User, got %v", decoded.Payload["name"])
	}
	if decoded.Header["alg"] != "HS256" {
		t.Errorf("expected alg=HS256, got %v", decoded.Header["alg"])
	}
}

// writeKeyToTempFile writes PEM data to a temp file and returns the path.
func writeKeyToTempFile(t *testing.T, pemData []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "key.pem")
	if err := os.WriteFile(path, pemData, 0o600); err != nil {
		t.Fatalf("failed to write temp key file: %v", err)
	}
	return path
}

func TestSignRSA_Generated(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der := x509.MarshalPKCS1PrivateKey(rsaKey)
	pemData := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})
	keyPath := writeKeyToTempFile(t, pemData)

	sig, err := signRSA("RS256", keyPath, "test.payload")
	if err != nil {
		t.Fatalf("signRSA unexpected error: %v", err)
	}
	if len(sig) == 0 {
		t.Error("expected non-empty signature")
	}
}

func TestSignRSA_MissingKeyFile(t *testing.T) {
	_, err := signRSA("RS256", "", "test.payload")
	if err == nil {
		t.Fatal("expected error for empty key file")
	}
}

func TestSignRSA_NonexistentFile(t *testing.T) {
	_, err := signRSA("RS256", "/tmp/nonexistent-key-file-12345.pem", "test.payload")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestSignECDSA_Generated(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
	keyPath := writeKeyToTempFile(t, pemData)

	sig, err := signECDSA("ES256", keyPath, "test.payload")
	if err != nil {
		t.Fatalf("signECDSA unexpected error: %v", err)
	}
	// ES256 produces a 64-byte signature (32 bytes r + 32 bytes s)
	if len(sig) != 64 {
		t.Errorf("expected 64-byte signature for ES256, got %d", len(sig))
	}
}

func TestSignECDSA_CurveMismatch(t *testing.T) {
	// Generate P-384 key but try to sign with ES256 (expects P-256)
	ecKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
	keyPath := writeKeyToTempFile(t, pemData)

	_, err = signECDSA("ES256", keyPath, "test.payload")
	if err == nil {
		t.Fatal("expected error for curve mismatch")
	}
}

func TestSignECDSA_MissingKeyFile(t *testing.T) {
	_, err := signECDSA("ES256", "", "test.payload")
	if err == nil {
		t.Fatal("expected error for empty key file")
	}
}

func TestSign_RS256_EndToEnd(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	der := x509.MarshalPKCS1PrivateKey(rsaKey)
	pemData := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der})
	keyPath := writeKeyToTempFile(t, pemData)

	token, err := Sign(SignOptions{
		Algorithm: "RS256",
		KeyFile:   keyPath,
		Claims:    `{"sub":"rs256-test","iss":"jwx"}`,
	})
	if err != nil {
		t.Fatalf("Sign(RS256) unexpected error: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3-part token, got %d parts", len(parts))
	}

	decoded, err := Decode(token)
	if err != nil {
		t.Fatalf("Decode() unexpected error: %v", err)
	}
	if decoded.Header["alg"] != "RS256" {
		t.Errorf("expected alg=RS256, got %v", decoded.Header["alg"])
	}
	if decoded.Payload["sub"] != "rs256-test" {
		t.Errorf("expected sub=rs256-test, got %v", decoded.Payload["sub"])
	}
	if decoded.Payload["iss"] != "jwx" {
		t.Errorf("expected iss=jwx, got %v", decoded.Payload["iss"])
	}
}

func TestIsSupported(t *testing.T) {
	valid := []string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512", "ES256", "ES384", "ES512", "EdDSA"}
	for _, alg := range valid {
		if !isSupported(alg) {
			t.Errorf("expected %s to be supported", alg)
		}
	}

	// Case-insensitive check
	if !isSupported("hs256") {
		t.Error("expected hs256 (lowercase) to be supported")
	}
	if !isSupported("eddsa") {
		t.Error("expected eddsa (lowercase) to be supported")
	}

	invalid := []string{"FAKE", "RS128", "HS1024", "none", ""}
	for _, alg := range invalid {
		if isSupported(alg) {
			t.Errorf("expected %s to NOT be supported", alg)
		}
	}
}

func TestSign_ExtraHeaderFields(t *testing.T) {
	token, err := Sign(SignOptions{
		Algorithm: "HS256",
		Secret:    "test-secret",
		Claims:    `{"sub":"header-test"}`,
		Header: map[string]any{
			"kid": "my-key-id",
			"jku": "https://example.com/jwks",
		},
	})
	if err != nil {
		t.Fatalf("Sign() with extra headers unexpected error: %v", err)
	}

	decoded, err := Decode(token)
	if err != nil {
		t.Fatalf("Decode() unexpected error: %v", err)
	}

	if decoded.Header["kid"] != "my-key-id" {
		t.Errorf("expected kid=my-key-id, got %v", decoded.Header["kid"])
	}
	if decoded.Header["jku"] != "https://example.com/jwks" {
		t.Errorf("expected jku=https://example.com/jwks, got %v", decoded.Header["jku"])
	}
	if decoded.Header["alg"] != "HS256" {
		t.Errorf("expected alg=HS256, got %v", decoded.Header["alg"])
	}
	if decoded.Header["typ"] != "JWT" {
		t.Errorf("expected typ=JWT, got %v", decoded.Header["typ"])
	}
}

func TestSign_ES256_EndToEnd(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate EC key: %v", err)
	}
	der, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		t.Fatalf("failed to marshal EC key: %v", err)
	}
	pemData := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
	keyPath := writeKeyToTempFile(t, pemData)

	token, err := Sign(SignOptions{
		Algorithm: "ES256",
		KeyFile:   keyPath,
		Claims:    `{"sub":"ec-test"}`,
	})
	if err != nil {
		t.Fatalf("Sign(ES256) unexpected error: %v", err)
	}

	decoded, err := Decode(token)
	if err != nil {
		t.Fatalf("Decode() unexpected error: %v", err)
	}
	if decoded.Header["alg"] != "ES256" {
		t.Errorf("expected alg=ES256, got %v", decoded.Header["alg"])
	}
	if decoded.Payload["sub"] != "ec-test" {
		t.Errorf("expected sub=ec-test, got %v", decoded.Payload["sub"])
	}
}

func TestSign_CaseInsensitiveAlgorithm(t *testing.T) {
	token, err := Sign(SignOptions{
		Algorithm: "hs256",
		Secret:    "secret",
		Claims:    `{"sub":"case"}`,
	})
	if err != nil {
		t.Fatalf("Sign() with lowercase alg unexpected error: %v", err)
	}
	decoded, err := Decode(token)
	if err != nil {
		t.Fatalf("Decode() unexpected error: %v", err)
	}
	if decoded.Header["alg"] != "HS256" {
		t.Errorf("expected alg=HS256 (uppercased), got %v", decoded.Header["alg"])
	}
}
