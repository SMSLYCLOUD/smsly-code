package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SMSLYUser represents a user authenticated via SMSLY Identity.
type SMSLYUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

// Config holds the configuration for SMSLY OAuth provider.
type Config struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	CallbackURL  string
}

// SMSLYProvider handles authentication with SMSLY Identity.
type SMSLYProvider struct {
	config Config
	client *http.Client
}

// NewSMSLYProvider creates a new provider instance using environment variables.
func NewSMSLYProvider() (*SMSLYProvider, error) {
	config := Config{
		ClientID:     os.Getenv("SMSLY_IDENTITY_CLIENT_ID"),
		ClientSecret: os.Getenv("SMSLY_IDENTITY_CLIENT_SECRET"),
		AuthURL:      os.Getenv("SMSLY_IDENTITY_AUTH_URL"),
		TokenURL:     os.Getenv("SMSLY_IDENTITY_TOKEN_URL"),
		UserInfoURL:  os.Getenv("SMSLY_IDENTITY_USERINFO_URL"),
		CallbackURL:  os.Getenv("SMSLY_IDENTITY_CALLBACK_URL"),
	}

	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, errors.New("SMSLY_IDENTITY_CLIENT_ID and SMSLY_IDENTITY_CLIENT_SECRET are required")
	}

	// Set defaults for URLs if not provided (useful for dev/test without full env)
	if config.AuthURL == "" {
		config.AuthURL = "https://identity.smsly.cloud/oauth/authorize"
	}
	if config.TokenURL == "" {
		config.TokenURL = "https://identity.smsly.cloud/oauth/token"
	}
	if config.UserInfoURL == "" {
		config.UserInfoURL = "https://identity.smsly.cloud/api/v1/user"
	}

	return &SMSLYProvider{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// GenerateAuthURL returns the URL to redirect the user for authentication.
func (p *SMSLYProvider) GenerateAuthURL(state string) string {
	u, _ := url.Parse(p.config.AuthURL)
	q := u.Query()
	q.Set("client_id", p.config.ClientID)
	q.Set("redirect_uri", p.config.CallbackURL)
	q.Set("response_type", "code")
	q.Set("state", state)
	q.Set("scope", "openid profile email")
	u.RawQuery = q.Encode()
	return u.String()
}

// Exchange exchanges an authorization code for a user profile.
func (p *SMSLYProvider) Exchange(ctx context.Context, code string) (*SMSLYUser, error) {
	// 1. Exchange code for access token
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.CallbackURL)
	data.Set("client_id", p.config.ClientID)
	data.Set("client_secret", p.config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", p.config.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// 2. Use access token to fetch user info
	req, err = http.NewRequestWithContext(ctx, "GET", p.config.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Accept", "application/json")

	resp, err = p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute user info request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info fetch failed (status %d): %s", resp.StatusCode, string(body))
	}

	var user SMSLYUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &user, nil
}
