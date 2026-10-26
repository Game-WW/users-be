package session

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Game-WW/users-be/pkg/errs"
)

var ErrExpired = fmt.Errorf("session expired: %w", errs.ErrUnauthenticated)

type State struct {
	UserID    string
	TokenHash string
	UserAgent string
	IP        string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Session struct {
	state State
}

func New(userID, userAgent, ip string, now time.Time, ttl time.Duration) (*Session, string) {
	token := rand.Text()

	return &Session{state: State{
		UserID:    userID,
		TokenHash: HashToken(token),
		UserAgent: userAgent,
		IP:        ip,
		CreatedAt: now,
		ExpiresAt: now.Add(ttl),
	}}, token
}

func Restore(state State) *Session {
	return &Session{state: state}
}

func (s *Session) UserID() string {
	return s.state.UserID
}

func (s *Session) Snapshot() State {
	return s.state
}

func (s *Session) Rotate(now time.Time, ttl time.Duration) (string, error) {
	if !now.Before(s.state.ExpiresAt) {
		return "", ErrExpired
	}

	token := rand.Text()
	s.state.TokenHash = HashToken(token)
	s.state.ExpiresAt = now.Add(ttl)

	return token, nil
}
