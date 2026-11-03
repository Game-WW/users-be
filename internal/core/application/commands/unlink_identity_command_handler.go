package commands

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/user"
)

type IdentityUnlicker interface {
	UnlinkIdentity(ctx context.Context, userId string, provider user.Provider) error
}

type UnlinkIdentityHandler struct {
	repo IdentityUnlicker
}

func NewUnlinkIdentityHandler(repo IdentityUnlicker) *UnlinkIdentityHandler {
	return &UnlinkIdentityHandler{repo: repo}
}

func (h *UnlinkIdentityHandler) Handle(ctx context.Context, cmd UnlinkIdentityCommand) error {
	if err := h.repo.UnlinkIdentity(ctx, cmd.UserID, cmd.Provider); err != nil {
		return fmt.Errorf("unlink identity: %w", err)
	}

	return nil
}
