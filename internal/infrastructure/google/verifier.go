package google

import (
	"context"
	"fmt"

	"google.golang.org/api/idtoken"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

var issuers = map[string]bool{
	"accounts.google.com":         true,
	"https://accounts.google.com": true,
}

type Verifier struct {
	clientID string
}

func NewVerifier(clientID string) *Verifier {
	return &Verifier{clientID: clientID}
}

func (v *Verifier) Verify(ctx context.Context, idToken string) (user.ExternalIdentity, error) {
	payload, err := idtoken.Validate(ctx, idToken, v.clientID)
	if err != nil {
		return user.ExternalIdentity{}, fmt.Errorf("validate google id_token: %w: %w", errs.ErrUnauthenticated, err)
	}
	if !issuers[payload.Issuer] {
		return user.ExternalIdentity{}, fmt.Errorf("google id_token issuer %q: %w", payload.Issuer, errs.ErrUnauthenticated)
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	return user.ExternalIdentity{
		Provider: user.ProviderGoogle,
		Subject:  payload.Subject,
		Email:    email,
		Name:     name,
	}, nil
}
