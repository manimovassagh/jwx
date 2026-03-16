package jwt

import (
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
