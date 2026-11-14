package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

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
	transport := &http.Transport{
		MaxConnsPerHost:     10,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	return &Verifier{
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		},
	}
}

func (v *Verifier) Verify(ctx context.Context, creds commands.AuthCredentials) (user.ExternalIdentity, error) {
	if creds.Token != "" && creds.Code != "" {
		return user.ExternalIdentity{}, fmt.Errorf("ambiguous yandex credentials: %w", errs.ErrInvalidArgument)
	}
	accessToken := creds.Token
	if accessToken == "" {
		if creds.Code == "" {
			return user.ExternalIdentity{}, fmt.Errorf("empty yandex credentials: %w", errs.ErrUnauthenticated)
		}
		var err error
		accessToken, err = v.exchangeCode(ctx, creds.Code, creds.CodeVerifier)
		if err != nil {
			return user.ExternalIdentity{}, err
		}
	}

	ydUser, err := v.getUserInfo(ctx, accessToken)
	if err != nil {
		return user.ExternalIdentity{}, err
	}

	if ydUser.ID == "" || (creds.Token != "" && ydUser.ClientID != v.clientID) {
		return user.ExternalIdentity{}, fmt.Errorf("invalid yandex identity or token client: %w", errs.ErrUnauthenticated)
	}

	email := strings.TrimSpace(ydUser.DefaultEmail)
	if email == "" {
		for _, candidate := range ydUser.Emails {
			if candidate = strings.TrimSpace(candidate); candidate != "" {
				email = candidate
				break
			}
		}
	}

	return user.ExternalIdentity{
		Provider: user.ProviderYandex,
		Subject:  ydUser.ID,
		Email:    email,
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read token response body: %w", err)
	}
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
	ClientID     string   `json:"client_id"`
	ID           string   `json:"id"`
	DefaultEmail string   `json:"default_email"`
	Emails       []string `json:"emails"`
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read user info response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("yandex user info status %d: %w", resp.StatusCode, errs.ErrUnauthenticated)
		}
		return nil, fmt.Errorf("yandex user info status %d", resp.StatusCode)
	}

	var u ydUser
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}
