// Package identity выбирает проверку id_token по провайдеру.
package identity

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

// ProviderVerifier — проверка id_token одного провайдера (google.Verifier, потом Apple и т.д.).
type ProviderVerifier interface {
	Verify(ctx context.Context, idToken string) (user.ExternalIdentity, error)
}

// Verifier — провайдер → его проверка. Новый провайдер добавляется в composition root.
type Verifier map[user.Provider]ProviderVerifier

var _ commands.IdentityVerifier = Verifier{}

func (v Verifier) Verify(ctx context.Context, provider user.Provider, idToken string) (user.ExternalIdentity, error) {
	pv, ok := v[provider]
	if !ok {
		// Провайдер известен, но выключен (не задан его client id) — для клиента это недопустимый провайдер, 400.
		return user.ExternalIdentity{}, fmt.Errorf("provider %q is disabled: %w", provider, errs.ErrInvalidArgument)
	}

	return pv.Verify(ctx, idToken)
}
