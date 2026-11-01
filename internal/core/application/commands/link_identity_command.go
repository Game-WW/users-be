package commands

import "github.com/Game-WW/users-be/internal/core/domain/user"

type LinkIdentityCommand struct {
	UserID   string
	Provider user.Provider
	Creds    AuthCredentials
}

type LinkIdentityResult struct {
	Provider user.Provider
	Email    string
}
