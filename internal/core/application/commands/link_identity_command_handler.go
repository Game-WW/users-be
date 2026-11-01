package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

type IdentityLinker interface {
	FindByIdentity(ctx context.Context, provider user.Provider, subject string) (*user.User, error)
	LinkIdentity(ctx context.Context, userID string, identity user.ExternalIdentity, now time.Time) error
}

type LinkIdentityHandler struct {
	identities map[user.Provider]IdentityVerifier
	repo       IdentityLinker
	clock      Clock
}

func NewLinkIdentityHandler(
	identities map[user.Provider]IdentityVerifier,
	repo IdentityLinker,
	clock Clock,
) *LinkIdentityHandler {
	return &LinkIdentityHandler{
		identities: identities,
		repo:       repo,
		clock:      clock,
	}
}

func (h *LinkIdentityHandler) Handle(ctx context.Context, cmd LinkIdentityCommand) (LinkIdentityResult, error) {
	verifier, ok := h.identities[cmd.Provider]
	if !ok {
		return LinkIdentityResult{}, fmt.Errorf("provider %s is not configured: %w", cmd.Provider, errs.ErrInvalidArgument)
	}

	identity, err := verifier.Verify(ctx, cmd.Creds)
	if err != nil {
		return LinkIdentityResult{}, fmt.Errorf("verify identity: %w", err)
	}

	existingUser, err := h.repo.FindByIdentity(ctx, identity.Provider, identity.Subject)
	if err == nil {
		if existingUser.ID() == cmd.UserID {
			return LinkIdentityResult{Provider: identity.Provider, Email: identity.Email}, nil
		}

		return LinkIdentityResult{}, fmt.Errorf("identity is already linked to another user: %w", ErrIdentityTaken)
	}

	if !errors.Is(err, errs.ErrNotFound) {
		return LinkIdentityResult{}, fmt.Errorf("check identity: %w", err)
	}

	err = h.repo.LinkIdentity(ctx, cmd.UserID, identity, h.clock.Now())
	if err != nil {
		return LinkIdentityResult{}, fmt.Errorf("link identity: %w", err)
	}

	return LinkIdentityResult{
		Provider: identity.Provider,
		Email:    identity.Email,
	}, nil
}
