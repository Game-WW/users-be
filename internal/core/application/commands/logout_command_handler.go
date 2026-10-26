package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
)

type SessionRevoker interface {
	DeleteByTokenHash(ctx context.Context, hash string) error
}

type LogoutHandler struct {
	sessions SessionRevoker
}

func NewLogoutHandler(sessions SessionRevoker) *LogoutHandler {
	return &LogoutHandler{sessions: sessions}
}

func (h *LogoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	if cmd.RefreshToken == "" {
		return nil
	}

	err := h.sessions.DeleteByTokenHash(ctx, session.HashToken(cmd.RefreshToken))
	if err != nil && !errors.Is(err, errs.ErrNotFound) {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
