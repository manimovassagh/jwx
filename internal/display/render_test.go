package display

import (
	"strings"
	"testing"
	"time"

	"github.com/manimovassagh/jwx/internal/jwt"
)

func makeBasicToken() *jwt.DecodedToken {
	return &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"sub":  "1234567890",
			"name": "John Doe",
		},
		Signature: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk",
	}
}

func TestRender_BasicToken(t *testing.T) {
	token := makeBasicToken()
	output := Render(token)

	for _, want := range []string{"Header", "Payload", "Signature"} {
		if !strings.Contains(output, want) {
			t.Errorf("Render output missing %q section", want)
		}
	}
}

func TestRender_ExpiredToken(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour)
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"sub": "user1",
			"exp": float64(past.Unix()),
		},
		Signature: "abc123",
		IsExpired: true,
		ExpiresAt: &past,
	}

	output := Render(token)
	if !strings.Contains(output, "EXPIRED") {
		t.Error("Render output should contain 'EXPIRED' for an expired token")
	}
}

func TestRender_ValidFutureExpiry(t *testing.T) {
	future := time.Now().Add(24 * time.Hour)
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"sub": "user1",
			"exp": float64(future.Unix()),
		},
		Signature: "abc123",
		IsExpired: false,
		ExpiresAt: &future,
	}

	output := Render(token)
	if !strings.Contains(output, "Expires") {
		t.Error("Render output should contain 'Expires' for a valid future expiry")
	}
}

func TestRender_AllStandardClaims(t *testing.T) {
	now := time.Now()
	future := now.Add(1 * time.Hour)
	past := now.Add(-1 * time.Hour)

	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "RS256",
			"typ": "JWT",
		},
		Payload: map[string]interface{}{
			"iss": "https://example.com",
			"sub": "user123",
			"aud": "my-app",
			"exp": float64(future.Unix()),
			"iat": float64(now.Unix()),
			"nbf": float64(past.Unix()),
			"jti": "unique-id-abc",
		},
		Signature: "longsignaturevaluehere1234567890abcdef",
		IsExpired: false,
		ExpiresAt: &future,
		IssuedAt:  &now,
		NotBefore: &past,
	}

	output := Render(token)

	for _, claim := range []string{"iss", "sub", "aud", "exp", "iat", "nbf", "jti"} {
		if !strings.Contains(output, claim+":") {
			t.Errorf("Render output missing standard claim %q", claim)
		}
	}
}

func TestRenderClaims_ValueTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"string value", "hello", `"hello"`},
		{"integer number", float64(42), "42"},
		{"float number", float64(3.14), "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"nil value", nil, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := map[string]interface{}{
				"key": tt.value,
			}
			output := renderClaims(claims)
			if !strings.Contains(output, tt.expected) {
				t.Errorf("renderClaims(%v) = %q, want it to contain %q", tt.value, output, tt.expected)
			}
		})
	}
}

func TestOrderKeys_Prioritization(t *testing.T) {
	keys := []string{"custom", "aud", "exp", "iss", "name", "sub"}
	priority := []string{"iss", "sub", "aud", "exp", "nbf", "iat", "jti"}

	result := orderKeys(keys, priority)

	// Priority keys that exist should come first, in priority order
	expectedPrefix := []string{"iss", "sub", "aud", "exp"}
	for i, want := range expectedPrefix {
		if i >= len(result) {
			t.Fatalf("result too short: got %d elements, want at least %d", len(result), i+1)
		}
		if result[i] != want {
			t.Errorf("orderKeys result[%d] = %q, want %q", i, result[i], want)
		}
	}

	// Remaining non-priority keys should follow
	remaining := result[len(expectedPrefix):]
	for _, r := range remaining {
		for _, p := range expectedPrefix {
			if r == p {
				t.Errorf("priority key %q found in remaining section", r)
			}
		}
	}

	// Total length should match
	if len(result) != len(keys) {
		t.Errorf("orderKeys returned %d keys, want %d", len(result), len(keys))
	}
}

func TestOrderKeys_NoPriorityKeys(t *testing.T) {
	keys := []string{"z_custom", "a_custom", "m_custom"}
	priority := []string{"iss", "sub"}

	result := orderKeys(keys, priority)

	if len(result) != len(keys) {
		t.Fatalf("expected %d keys, got %d", len(keys), len(result))
	}

	// Should just be the original keys in order (since none match priority)
	for i, k := range keys {
		if result[i] != k {
			t.Errorf("result[%d] = %q, want %q", i, result[i], k)
		}
	}
}

func TestOrderKeys_AllPriorityKeys(t *testing.T) {
	keys := []string{"exp", "iss", "sub"}
	priority := []string{"iss", "sub", "exp"}

	result := orderKeys(keys, priority)

	expected := []string{"iss", "sub", "exp"}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("result[%d] = %q, want %q", i, result[i], want)
		}
	}
}

func TestOrderKeys_EmptyInputs(t *testing.T) {
	result := orderKeys(nil, nil)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}

	result = orderKeys([]string{"a", "b"}, nil)
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestRender_SignatureTruncation(t *testing.T) {
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
		},
		Payload:   map[string]interface{}{},
		Signature: "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJ",
	}

	output := Render(token)
	if !strings.Contains(output, "...") {
		t.Error("long signature should be truncated with '...'")
	}
}

func TestRender_ShortSignatureNotTruncated(t *testing.T) {
	token := &jwt.DecodedToken{
		Header: map[string]interface{}{
			"alg": "HS256",
		},
		Payload:   map[string]interface{}{},
		Signature: "short",
	}

	output := Render(token)
	if !strings.Contains(output, "short") {
		t.Error("short signature should appear in output")
	}
}

func TestRender_NoColor(t *testing.T) {
	// Enable NoColor mode
	NoColor = true
	defer func() {
		NoColor = false
		initStyles()
	}()

	token := makeBasicToken()
	output := Render(token)

	// ANSI escape codes start with ESC (\x1b)
	if strings.Contains(output, "\x1b[") {
		t.Errorf("NoColor output should not contain ANSI escape codes, got:\n%s", output)
	}

	// Content should still be present
	for _, want := range []string{"Header", "Payload", "Signature", "HS256"} {
		if !strings.Contains(output, want) {
			t.Errorf("NoColor output missing %q", want)
		}
	}
}

func TestRender_UnknownAlgorithm(t *testing.T) {
	token := &jwt.DecodedToken{
		Header:    map[string]interface{}{},
		Payload:   map[string]interface{}{},
		Signature: "sig",
	}

	output := Render(token)
	if !strings.Contains(output, "unknown") {
		t.Error("missing alg should render as 'unknown'")
	}
}
