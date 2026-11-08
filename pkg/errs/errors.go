package errs

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrNicknameTaken   = errors.New("nickname already taken")
	ErrLastIdentity    = errors.New("cannot unlink the last identity")
)
