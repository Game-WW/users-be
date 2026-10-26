package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
)

type SessionReader interface {
	FindByTokenHash(ctx context.Context, hash string) (*session.Session, error)
}

type RefreshHandler struct {
	sessions  SessionReader
	access    AccessTokenIssuer
	clock     Clock
	accessTTL time.Duration
}

func NewRefreshHandler(sessions SessionReader, access AccessTokenIssuer, clock Clock, accessTTL time.Duration) *RefreshHandler {
	return &RefreshHandler{sessions: sessions, access: access, clock: clock, accessTTL: accessTTL}
}

func (h *RefreshHandler) Handle(ctx context.Context, cmd RefreshCommand) (RefreshResult, error) {
	if cmd.RefreshToken == "" {
		return RefreshResult{}, errs.ErrUnauthenticated
	}

	sess, err := h.sessions.FindByTokenHash(ctx, session.HashToken(cmd.RefreshToken))
	if errors.Is(err, errs.ErrNotFound) {
		return RefreshResult{}, fmt.Errorf("refresh session not found: %w", errs.ErrUnauthenticated)
	}
	if err != nil {
		return RefreshResult{}, fmt.Errorf("find session: %w", err)
	}

	now := h.clock.Now()
	if err := sess.ValidateRefresh(now); err != nil {
		return RefreshResult{}, err
	}

	token, err := h.access.Issue(sess.UserID(), now.Add(h.accessTTL))
	if err != nil {
		return RefreshResult{}, fmt.Errorf("issue access token: %w", err)
	}

	return RefreshResult{AccessToken: token}, nil
}
