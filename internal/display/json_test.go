package display

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/manimovassagh/jwx/internal/jwt"
)

func TestRenderJSON_ProducesValidJSON(t *testing.T) {
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"sub":  "1234567890",
			"name": "John Doe",
		},
		Signature: "test-signature",
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	if !json.Valid([]byte(output)) {
		t.Errorf("RenderJSON output is not valid JSON: %s", output)
	}
}

func TestRenderJSON_HasCorrectFields(t *testing.T) {
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "RS256",
		},
		Payload: map[string]interface{}{
			"sub": "user1",
		},
		Signature: "sig-data",
		IsExpired: false,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	requiredFields := []string{"header", "payload", "signature", "is_expired"}
	for _, field := range requiredFields {
		if _, ok := result[field]; !ok {
			t.Errorf("RenderJSON output missing required field %q", field)
		}
	}
}

func TestRenderJSON_ExpiredTokenSetsIsExpired(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour)
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
		},
		Payload: map[string]interface{}{
			"exp": float64(past.Unix()),
		},
		Signature: "sig",
		IsExpired: true,
		ExpiresAt: &past,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if !result.IsExpired {
		t.Error("RenderJSON should set is_expired=true for an expired token")
	}
	if result.ExpiresAt == nil {
		t.Error("expected expires_at to be set for token with ExpiresAt")
	}
}

func TestRenderJSON_NotExpiredToken(t *testing.T) {
	future := time.Now().Add(1 * time.Hour)
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
		},
		Payload: map[string]interface{}{
			"exp": float64(future.Unix()),
		},
		Signature: "sig",
		IsExpired: false,
		ExpiresAt: &future,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if result.IsExpired {
		t.Error("RenderJSON should set is_expired=false for a valid token")
	}
	if result.ExpiresAt == nil {
		t.Error("expected expires_at to be set for token with ExpiresAt")
	}
}

func TestRenderJSON_Roundtrip(t *testing.T) {
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "ES256",
			"typ": "JWT",
			"kid": "key-1",
		},
		Payload: map[string]interface{}{
			"iss":   "https://auth.example.com",
			"sub":   "user-42",
			"aud":   "my-client",
			"exp":   float64(1700000000),
			"iat":   float64(1699999000),
			"admin": true,
			"score": float64(9.5),
		},
		Signature: "roundtrip-signature-data",
		IsExpired: true,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	// Check header values
	if alg, ok := result.Header["alg"].(string); !ok || alg != "ES256" {
		t.Errorf("header.alg = %v, want %q", result.Header["alg"], "ES256")
	}
	if typ, ok := result.Header["typ"].(string); !ok || typ != "JWT" {
		t.Errorf("header.typ = %v, want %q", result.Header["typ"], "JWT")
	}
	if kid, ok := result.Header["kid"].(string); !ok || kid != "key-1" {
		t.Errorf("header.kid = %v, want %q", result.Header["kid"], "key-1")
	}

	// Check payload values
	if iss, ok := result.Payload["iss"].(string); !ok || iss != "https://auth.example.com" {
		t.Errorf("payload.iss = %v, want %q", result.Payload["iss"], "https://auth.example.com")
	}
	if sub, ok := result.Payload["sub"].(string); !ok || sub != "user-42" {
		t.Errorf("payload.sub = %v, want %q", result.Payload["sub"], "user-42")
	}
	if admin, ok := result.Payload["admin"].(bool); !ok || admin != true {
		t.Errorf("payload.admin = %v, want true", result.Payload["admin"])
	}

	// Check signature
	if result.Signature != "roundtrip-signature-data" {
		t.Errorf("signature = %q, want %q", result.Signature, "roundtrip-signature-data")
	}

	// Check is_expired
	if !result.IsExpired {
		t.Error("is_expired should be true")
	}
}

func TestRenderJSON_TimestampFields(t *testing.T) {
	exp := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
	iat := time.Date(2025, 1, 15, 11, 0, 0, 0, time.UTC)
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
		},
		Payload: map[string]interface{}{
			"exp": float64(exp.Unix()),
			"iat": float64(iat.Unix()),
		},
		Signature: "sig",
		IsExpired: true,
		ExpiresAt: &exp,
		IssuedAt:  &iat,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if result.ExpiresAt == nil {
		t.Fatal("expected expires_at to be set")
	}
	if *result.ExpiresAt != "2025-01-15T12:00:00Z" {
		t.Errorf("expected expires_at=%q, got %q", "2025-01-15T12:00:00Z", *result.ExpiresAt)
	}

	if result.IssuedAt == nil {
		t.Fatal("expected issued_at to be set")
	}
	if *result.IssuedAt != "2025-01-15T11:00:00Z" {
		t.Errorf("expected issued_at=%q, got %q", "2025-01-15T11:00:00Z", *result.IssuedAt)
	}
}

func TestRenderJSON_NoTimestamps(t *testing.T) {
	token := &jwt.DecodedToken{
		Header:    map[string]interface{}{"alg": "HS256"},
		Payload:   map[string]interface{}{"sub": "test"},
		Signature: "sig",
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if result.ExpiresAt != nil {
		t.Error("expected expires_at to be nil when not set")
	}
	if result.IssuedAt != nil {
		t.Error("expected issued_at to be nil when not set")
	}
}

func TestRenderJSON_EmptyPayload(t *testing.T) {
	token := &jwt.DecodedToken{
		Header:    map[string]interface{}{},
		Payload:   map[string]interface{}{},
		Signature: "",
		IsExpired: false,
	}

	output, err := RenderJSON(token)
	if err != nil {
		t.Fatalf("RenderJSON returned error: %v", err)
	}

	var result JSONOutput
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(result.Header) != 0 {
		t.Errorf("expected empty header, got %v", result.Header)
	}
	if len(result.Payload) != 0 {
		t.Errorf("expected empty payload, got %v", result.Payload)
	}
}
