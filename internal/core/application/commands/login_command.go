package commands

import "github.com/Game-WW/users-be/internal/core/domain/user"

type LoginCommand struct {
	Provider  user.Provider
	Creds     AuthCredentials
	UserAgent string
	IP        string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}
