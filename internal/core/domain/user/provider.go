package user

import (
	"fmt"

	"github.com/Game-WW/users-be/pkg/errs"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
)

func ParseProvider(s string) (Provider, error) {
	switch p := Provider(s); p {
	case ProviderGoogle:
		return p, nil
	default:
		return "", fmt.Errorf("unknown provider %q: %w", s, errs.ErrInvalidArgument)
	}
}
