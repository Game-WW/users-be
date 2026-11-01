package github

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

	ghUser, err := v.getUser(ctx, accessToken)
	if err != nil {
		return user.ExternalIdentity{}, err
	}

	email, err := v.getPrimaryEmail(ctx, accessToken)
	if err != nil {
		return user.ExternalIdentity{}, err
	}

	return user.ExternalIdentity{
		Provider: user.ProviderGitHub,
		Subject:  fmt.Sprintf("%d", ghUser.ID),
		Email:    email,
	}, nil
}

func (v *Verifier) exchangeCode(ctx context.Context, code, codeVerifier string) (string, error) {
	data := url.Values{}
	data.Set("client_id", v.clientID)
	data.Set("client_secret", v.clientSecret)
	data.Set("code", code)
	if codeVerifier != "" {
		data.Set("code_verifier", codeVerifier)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tokenResp struct {
		AccessToken      string `json:"access_token"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	decodeErr := json.Unmarshal(body, &tokenResp)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github auth error: %s (%s): %w", tokenResp.Error, tokenResp.ErrorDescription, errs.ErrUnauthenticated)
	}
	if decodeErr != nil {
		return "", fmt.Errorf("unmarshal token: %w", decodeErr)
	}
	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("github error: %s (%s): %w", tokenResp.Error, tokenResp.ErrorDescription, errs.ErrUnauthenticated)
	}

	return tokenResp.AccessToken, nil
}

type ghUser struct {
	ID int64 `json:"id"`
}

func (v *Verifier) getUser(ctx context.Context, accessToken string) (*ghUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github user error: %s, body: %s", resp.Status, string(body))
	}

	var u ghUser
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}

func (v *Verifier) getPrimaryEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github emails error: %s, body: %s", resp.Status, string(body))
	}

	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", fmt.Errorf("unmarshal emails: %w", err)
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", fmt.Errorf("no primary verified email found: %w", errs.ErrInvalidArgument)
}
