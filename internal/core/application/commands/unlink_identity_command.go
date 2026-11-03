package commands

import "github.com/Game-WW/users-be/internal/core/domain/user"

type UnlinkIdentityCommand struct {
	UserID   string
	Provider user.Provider
}
