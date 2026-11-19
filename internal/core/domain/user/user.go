package user

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Game-WW/users-be/pkg/errs"
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

const maxProfileFieldLength = 20

func Validate(nickname, location *string) error {
	var trimmedNickname, trimmedLocation string
	if nickname != nil {
		trimmedNickname = strings.TrimSpace(*nickname)
		if trimmedNickname == "" {
			return fmt.Errorf("nickname cannot be empty: %w", errs.ErrInvalidArgument)
		}
		if utf8.RuneCountInString(trimmedNickname) > maxProfileFieldLength {
			return fmt.Errorf("nickname exceeds %d characters: %w", maxProfileFieldLength, errs.ErrInvalidArgument)
		}
	}
	if location != nil {
		trimmedLocation = strings.TrimSpace(*location)
		if utf8.RuneCountInString(trimmedLocation) > maxProfileFieldLength {
			return fmt.Errorf("location exceeds %d characters: %w", maxProfileFieldLength, errs.ErrInvalidArgument)
		}
	}

	if nickname != nil {
		*nickname = trimmedNickname
	}
	if location != nil {
		*location = trimmedLocation
	}
	return nil
}
