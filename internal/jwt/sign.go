package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"strings"
)

// SupportedAlgorithms lists all supported signing algorithms.
var SupportedAlgorithms = []string{
	"HS256", "HS384", "HS512",
	"RS256", "RS384", "RS512",
	"ES256", "ES384", "ES512",
	"EdDSA",
}

// SignOptions holds the parameters for signing a JWT.
type SignOptions struct {
	Algorithm string
	Secret    string         // For HMAC algorithms
	KeyFile   string         // For RSA/EC/EdDSA algorithms
	Claims    string         // JSON string of claims
	Header    map[string]any // Extra header fields
}

// Sign creates a signed JWT from the given options.
func Sign(opts SignOptions) (string, error) {
	alg := strings.ToUpper(opts.Algorithm)
	if !isSupported(alg) {
		return "", fmt.Errorf("unsupported algorithm: %s\nsupported: %s", alg, strings.Join(SupportedAlgorithms, ", "))
	}

	// Parse claims
	var claims map[string]any
	if err := json.Unmarshal([]byte(opts.Claims), &claims); err != nil {
		return "", fmt.Errorf("invalid claims JSON: %w", err)
	}

	// Build header
	header := map[string]any{
		"alg": alg,
		"typ": "JWT",
	}
	for k, v := range opts.Header {
		header[k] = v
	}

	// Encode header and payload
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to encode header: %w", err)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to encode claims: %w", err)
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)
	signingInput := headerB64 + "." + claimsB64

	// Sign
	var signature []byte
	switch {
	case strings.HasPrefix(alg, "HS"):
		signature, err = signHMAC(alg, opts.Secret, signingInput)
	case strings.HasPrefix(alg, "RS"):
		signature, err = signRSA(alg, opts.KeyFile, signingInput)
	case strings.HasPrefix(alg, "ES"):
		signature, err = signECDSA(alg, opts.KeyFile, signingInput)
	case alg == "EDDSA":
		signature, err = signEdDSA(opts.KeyFile, signingInput)
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", alg)
	}
	if err != nil {
		return "", err
	}

	sigB64 := base64.RawURLEncoding.EncodeToString(signature)
	return signingInput + "." + sigB64, nil
}

func signHMAC(alg, secret, input string) ([]byte, error) {
	if secret == "" {
		return nil, fmt.Errorf("--secret is required for %s", alg)
	}

	var h func() hash.Hash
	switch alg {
	case "HS256":
		h = sha256.New
	case "HS384":
		h = sha512.New384
	case "HS512":
		h = sha512.New
	default:
		return nil, fmt.Errorf("unsupported HMAC algorithm: %s", alg)
	}

	mac := hmac.New(h, []byte(secret))
	mac.Write([]byte(input))
	return mac.Sum(nil), nil
}

func signRSA(alg, keyFile, input string) ([]byte, error) {
	if keyFile == "" {
		return nil, fmt.Errorf("--key is required for %s", alg)
	}

	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := parseRSAPrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	var hashFunc func() hash.Hash
	var cryptoHash crypto.Hash
	switch alg {
	case "RS256":
		hashFunc = sha256.New
		cryptoHash = crypto.SHA256
	case "RS384":
		hashFunc = sha512.New384
		cryptoHash = crypto.SHA384
	case "RS512":
		hashFunc = sha512.New
		cryptoHash = crypto.SHA512
	}

	hasher := hashFunc()
	hasher.Write([]byte(input))
	hashed := hasher.Sum(nil)

	return rsa.SignPKCS1v15(rand.Reader, key, cryptoHash, hashed)
}

func signECDSA(alg, keyFile, input string) ([]byte, error) {
	if keyFile == "" {
		return nil, fmt.Errorf("--key is required for %s", alg)
	}

	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := parseECPrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	var expectedCurve elliptic.Curve
	var hashFunc func() hash.Hash
	switch alg {
	case "ES256":
		expectedCurve = elliptic.P256()
		hashFunc = sha256.New
	case "ES384":
		expectedCurve = elliptic.P384()
		hashFunc = sha512.New384
	case "ES512":
		expectedCurve = elliptic.P521()
		hashFunc = sha512.New
	}

	if key.Curve != expectedCurve {
		return nil, fmt.Errorf("key curve %s does not match algorithm %s", key.Curve.Params().Name, alg)
	}

	hasher := hashFunc()
	hasher.Write([]byte(input))
	hashed := hasher.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, key, hashed)
	if err != nil {
		return nil, err
	}

	// Encode r and s as fixed-size big-endian bytes
	keySize := (key.Curve.Params().BitSize + 7) / 8
	rBytes := r.Bytes()
	sBytes := s.Bytes()

	sig := make([]byte, keySize*2)
	copy(sig[keySize-len(rBytes):keySize], rBytes)
	copy(sig[keySize*2-len(sBytes):], sBytes)

	return sig, nil
}

func signEdDSA(keyFile, input string) ([]byte, error) {
	if keyFile == "" {
		return nil, fmt.Errorf("--key is required for EdDSA")
	}

	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	key, err := parseEdPrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Ed25519 private key: %w", err)
	}

	return ed25519.Sign(key, []byte(input)), nil
}

func isSupported(alg string) bool {
	for _, a := range SupportedAlgorithms {
		if strings.EqualFold(a, alg) {
			return true
		}
	}
	return false
}
