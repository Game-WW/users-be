package errs

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrLastIdentity    = errors.New("cannot unlink the last identity")
)
