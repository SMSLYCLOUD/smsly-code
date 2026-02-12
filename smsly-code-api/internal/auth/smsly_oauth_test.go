package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSMSLYProvider_MissingEnv(t *testing.T) {
	// Clear env vars
	os.Unsetenv("SMSLY_IDENTITY_CLIENT_ID")
	os.Unsetenv("SMSLY_IDENTITY_CLIENT_SECRET")

	provider, err := NewSMSLYProvider()
	assert.Error(t, err)
	assert.Nil(t, provider)
}

func TestNewSMSLYProvider_Success(t *testing.T) {
	os.Setenv("SMSLY_IDENTITY_CLIENT_ID", "test-client")
	os.Setenv("SMSLY_IDENTITY_CLIENT_SECRET", "test-secret")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_ID")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_SECRET")

	provider, err := NewSMSLYProvider()
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "test-client", provider.config.ClientID)
}

func TestGenerateAuthURL(t *testing.T) {
	os.Setenv("SMSLY_IDENTITY_CLIENT_ID", "test-client")
	os.Setenv("SMSLY_IDENTITY_CLIENT_SECRET", "test-secret")
	os.Setenv("SMSLY_IDENTITY_CALLBACK_URL", "http://localhost:8080/callback")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_ID")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_SECRET")
	defer os.Unsetenv("SMSLY_IDENTITY_CALLBACK_URL")

	provider, _ := NewSMSLYProvider()
	authURL := provider.GenerateAuthURL("random-state")

	assert.Contains(t, authURL, "client_id=test-client")
	assert.Contains(t, authURL, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback")
	assert.Contains(t, authURL, "state=random-state")
	assert.Contains(t, authURL, "scope=openid+profile+email")
}

func TestExchange_Success(t *testing.T) {
	// Mock SMSLY Identity Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "authorization_code", r.FormValue("grant_type"))
			assert.Equal(t, "valid-code", r.FormValue("code"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"access_token": "mock-access-token",
				"token_type":   "Bearer",
			})
			return
		}
		if r.URL.Path == "/api/v1/user" {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, "Bearer mock-access-token", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(SMSLYUser{
				ID:        123,
				Username:  "jules",
				Email:     "jules@smsly.cloud",
				FullName:  "Jules Agent",
				AvatarURL: "https://avatar.url",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	os.Setenv("SMSLY_IDENTITY_CLIENT_ID", "test-client")
	os.Setenv("SMSLY_IDENTITY_CLIENT_SECRET", "test-secret")
	os.Setenv("SMSLY_IDENTITY_TOKEN_URL", server.URL+"/oauth/token")
	os.Setenv("SMSLY_IDENTITY_USERINFO_URL", server.URL+"/api/v1/user")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_ID")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_SECRET")
	defer os.Unsetenv("SMSLY_IDENTITY_TOKEN_URL")
	defer os.Unsetenv("SMSLY_IDENTITY_USERINFO_URL")

	provider, err := NewSMSLYProvider()
	require.NoError(t, err)

    // Override client to allow localhost redirects if needed, but httptest server URL is fine
    // Actually, NewSMSLYProvider creates a new client. We don't expose it to replace transport.
    // However, httptest creates a local server that the default http.Client can reach.
    // The only issue is if NewSMSLYProvider sets a timeout. It does (10s). Should be fine.

	user, err := provider.Exchange(context.Background(), "valid-code")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(123), user.ID)
	assert.Equal(t, "jules", user.Username)
}

func TestExchange_TokenError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid code"))
			return
		}
	}))
	defer server.Close()

	os.Setenv("SMSLY_IDENTITY_CLIENT_ID", "test-client")
	os.Setenv("SMSLY_IDENTITY_CLIENT_SECRET", "test-secret")
	os.Setenv("SMSLY_IDENTITY_TOKEN_URL", server.URL+"/oauth/token")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_ID")
	defer os.Unsetenv("SMSLY_IDENTITY_CLIENT_SECRET")
	defer os.Unsetenv("SMSLY_IDENTITY_TOKEN_URL")

	provider, _ := NewSMSLYProvider()
	user, err := provider.Exchange(context.Background(), "invalid-code")
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "token exchange failed")
}
