package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/session"
	"github.com/Game-WW/users-be/pkg/errs"
)

type SessionRotator interface {
	FindByTokenHash(ctx context.Context, hash string) (*session.Session, error)
	Rotate(ctx context.Context, oldHash string, s *session.Session) error
}

type RefreshHandler struct {
	sessions   SessionRotator
	access     AccessTokenIssuer
	clock      Clock
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewRefreshHandler(sessions SessionRotator, access AccessTokenIssuer, clock Clock, accessTTL, refreshTTL time.Duration) *RefreshHandler {
	return &RefreshHandler{sessions: sessions, access: access, clock: clock, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (h *RefreshHandler) Handle(ctx context.Context, cmd RefreshCommand) (RefreshResult, error) {
	if cmd.RefreshToken == "" {
		return RefreshResult{}, errs.ErrUnauthenticated
	}

	oldHash := session.HashToken(cmd.RefreshToken)
	sess, err := h.sessions.FindByTokenHash(ctx, oldHash)
	if errors.Is(err, errs.ErrNotFound) {
		return RefreshResult{}, fmt.Errorf("refresh session not found: %w", errs.ErrUnauthenticated)
	}
	if err != nil {
		return RefreshResult{}, fmt.Errorf("find session: %w", err)
	}

	now := h.clock.Now()
	refreshToken, err := sess.Rotate(now, h.refreshTTL)
	if err != nil {
		return RefreshResult{}, err
	}

	err = h.sessions.Rotate(ctx, oldHash, sess)
	if errors.Is(err, errs.ErrNotFound) {
		return RefreshResult{}, fmt.Errorf("refresh token already rotated: %w", errs.ErrUnauthenticated)
	}
	if err != nil {
		return RefreshResult{}, fmt.Errorf("rotate session: %w", err)
	}

	accessToken, err := h.access.Issue(sess.UserID(), now.Add(h.accessTTL))
	if err != nil {
		return RefreshResult{}, fmt.Errorf("issue access token: %w", err)
	}

	return RefreshResult{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
