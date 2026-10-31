package queries

import (
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/user"
)

type ListIdentitiesQuery struct {
	UserID string
}

type Identity struct {
	Provider user.Provider
	Linked   bool
	Email    *string
	LinkedAt *time.Time
}

type ListIdentitiesResult struct {
	Identities []Identity
}
