package user

import (
	"fmt"
	"slices"

	"github.com/Game-WW/users-be/pkg/errs"
)

type Provider string

const (
	ProviderGoogle Provider = "google"
	ProviderGitHub Provider = "github"
	ProviderYandex Provider = "yandex"
)

var providers = []Provider{ProviderGoogle, ProviderGitHub, ProviderYandex}

func Providers() []Provider {
	return slices.Clone(providers)
}

func ParseProvider(s string) (Provider, error) {
	p := Provider(s)
	if !slices.Contains(providers, p) {
		return "", fmt.Errorf("unknown provider %q: %w", s, errs.ErrInvalidArgument)
	}

	return p, nil
}
