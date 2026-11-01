package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

type Verifier struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

func NewVerifier(clientID, clientSecret string) *Verifier {
	return &Verifier{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{},
	}
}

func (v *Verifier) Verify(ctx context.Context, creds commands.AuthCredentials) (user.ExternalIdentity, error) {
	code := creds.Code
	if code == "" {
		return user.ExternalIdentity{}, fmt.Errorf("empty code: %w", errs.ErrUnauthenticated)
	}

	accessToken, err := v.exchangeCode(ctx, code, creds.CodeVerifier)
	if err != nil {
		return user.ExternalIdentity{}, err
	}

	ydUser, err := v.getUserInfo(ctx, accessToken)
	if err != nil {
		return user.ExternalIdentity{}, err
	}

	return user.ExternalIdentity{
		Provider: user.ProviderYandex,
		Subject:  ydUser.ID,
		Email:    ydUser.DefaultEmail,
	}, nil
}

func (v *Verifier) exchangeCode(ctx context.Context, code, codeVerifier string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", v.clientID)
	data.Set("client_secret", v.clientSecret)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://oauth.yandex.ru/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("yandex auth error: %s, body: %s: %w", resp.Status, string(body), errs.ErrUnauthenticated)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("unmarshal token: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("yandex error: %s: %w", tokenResp.Error, errs.ErrUnauthenticated)
	}
	return tokenResp.AccessToken, nil
}

type ydUser struct {
	ID           string `json:"id"`
	DefaultEmail string `json:"default_email"`
}

func (v *Verifier) getUserInfo(ctx context.Context, accessToken string) (*ydUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "OAuth "+accessToken)

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex user info error: %s, body: %s", resp.Status, string(body))
	}

	var u ydUser
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}
