package avatar

import "errors"

var (
	ErrInvalidImage      = errors.New("invalid image file")
	ErrTooLarge          = errors.New("avatar file is too large")
	ErrUnsupportedFormat = errors.New("unsupported avatar format")
	ErrInvalidDimensions = errors.New("invalid avatar dimensions")
	ErrInvalidUserID     = errors.New("invalid user id")
)
