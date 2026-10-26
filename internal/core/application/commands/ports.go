package commands

import "time"

type Clock interface {
	Now() time.Time
}

type AccessTokenIssuer interface {
	Issue(userID string, expiresAt time.Time) (string, error)
}
