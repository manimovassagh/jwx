package jwt

import (
	"testing"
	"time"
)

// Test token: {"alg":"HS256","typ":"JWT"}.{"sub":"1234567890","name":"John Doe","iat":1516239022}.signature
const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"

// Expired token with exp=1516239022 (2018)
const expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZXhwIjoxNTE2MjM5MDIyfQ.4lKlMEaOmgSFB2ACmiTVsEL6L-5hAeTnjIPAhCMKBFE"

func TestDecode(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantErr   bool
		checkFunc func(t *testing.T, dt *DecodedToken)
	}{
		{
			name:    "valid token",
			token:   testToken,
			wantErr: false,
			checkFunc: func(t *testing.T, dt *DecodedToken) {
				if dt.Header["alg"] != "HS256" {
					t.Errorf("expected alg=HS256, got %v", dt.Header["alg"])
				}
				if dt.Header["typ"] != "JWT" {
					t.Errorf("expected typ=JWT, got %v", dt.Header["typ"])
				}
				if dt.Payload["sub"] != "1234567890" {
					t.Errorf("expected sub=1234567890, got %v", dt.Payload["sub"])
				}
				if dt.Payload["name"] != "John Doe" {
					t.Errorf("expected name=John Doe, got %v", dt.Payload["name"])
				}
			},
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: false,
			checkFunc: func(t *testing.T, dt *DecodedToken) {
				if !dt.IsExpired {
					t.Error("expected token to be expired")
				}
				if dt.ExpiresAt == nil {
					t.Fatal("expected ExpiresAt to be set")
				}
				expected := time.Unix(1516239022, 0)
				if !dt.ExpiresAt.Equal(expected) {
					t.Errorf("expected ExpiresAt=%v, got %v", expected, dt.ExpiresAt)
				}
			},
		},
		{
			name:    "empty string",
			token:   "",
			wantErr: true,
		},
		{
			name:    "not a jwt",
			token:   "this.is.not-valid-base64!",
			wantErr: true,
		},
		{
			name:    "only two parts",
			token:   "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0",
			wantErr: true,
		},
		{
			name:    "token with whitespace",
			token:   "  " + testToken + "  \n",
			wantErr: false,
			checkFunc: func(t *testing.T, dt *DecodedToken) {
				if dt.Header["alg"] != "HS256" {
					t.Errorf("expected alg=HS256, got %v", dt.Header["alg"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dt, err := Decode(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil && dt != nil {
				tt.checkFunc(t, dt)
			}
		})
	}
}
