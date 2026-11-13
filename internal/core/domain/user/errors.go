package user

import "errors"

var (
	ErrIdentityTaken   = errors.New("identity already linked to a user")
	ErrInvalidProvider = errors.New("invalid provider")
	ErrLastIdentity    = errors.New("cannot unlink the last identity")
	ErrNicknameTaken   = errors.New("nickname already taken")
)
