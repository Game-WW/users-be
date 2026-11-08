package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type State struct {
	ID        string
	Nickname  *string
	CreatedAt time.Time
}

type User struct {
	state State
}

func New(identity ExternalIdentity, now time.Time) (*User, error) {
	if identity.Subject == "" {
		return nil, errors.New("user: empty identity subject")
	}

	return &User{state: State{
		ID:        uuid.NewString(),
		CreatedAt: now,
	}}, nil
}

func Restore(state State) *User {
	return &User{state: state}
}

func (u *User) ID() string {
	return u.state.ID
}

func (u *User) Snapshot() State {
	return u.state
}
