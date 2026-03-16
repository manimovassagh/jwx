package display

import (
	"encoding/json"
	"time"

	"github.com/manimovassagh/jwx/internal/jwt"
)

// JSONOutput represents the structured JSON output of a decoded token.
type JSONOutput struct {
	Header    map[string]interface{} `json:"header"`
	Payload   map[string]interface{} `json:"payload"`
	Signature string                 `json:"signature"`
	IsExpired bool                   `json:"is_expired"`
	ExpiresAt *string                `json:"expires_at,omitempty"`
	IssuedAt  *string                `json:"issued_at,omitempty"`
}

// RenderJSON returns the token as formatted JSON for piping.
func RenderJSON(token *jwt.DecodedToken) (string, error) {
	out := JSONOutput{
		Header:    token.Header,
		Payload:   token.Payload,
		Signature: token.Signature,
		IsExpired: token.IsExpired,
	}

	if token.ExpiresAt != nil {
		s := token.ExpiresAt.Format(time.RFC3339)
		out.ExpiresAt = &s
	}
	if token.IssuedAt != nil {
		s := token.IssuedAt.Format(time.RFC3339)
		out.IssuedAt = &s
	}

	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
