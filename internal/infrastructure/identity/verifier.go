package identity

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

type ProviderVerifier interface {
	Verify(ctx context.Context, idToken string) (user.ExternalIdentity, error)
}

type Verifier map[user.Provider]ProviderVerifier

var _ commands.IdentityVerifier = Verifier{}

func (v Verifier) Verify(ctx context.Context, provider user.Provider, idToken string) (user.ExternalIdentity, error) {
	pv, ok := v[provider]
	if !ok {
		return user.ExternalIdentity{}, fmt.Errorf("provider %q is disabled: %w", provider, errs.ErrInvalidArgument)
	}

	return pv.Verify(ctx, idToken)
}
