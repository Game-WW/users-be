package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type State struct {
	ID          string
	UserName    *string
	Email       *string
	PhoneNumber *string
	CreatedAt   time.Time
}

type User struct {
	state State
}

func New(identity ExternalIdentity, now time.Time) (*User, error) {
	if identity.Subject == "" {
		return nil, errors.New("user: empty identity subject")
	}

	u := &User{state: State{ID: uuid.NewString(), CreatedAt: now}}
	if identity.Name != "" {
		u.state.UserName = &identity.Name
	}
	if identity.Email != "" {
		u.state.Email = &identity.Email
	}

	return u, nil
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
