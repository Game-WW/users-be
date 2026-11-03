package user

import "errors"

var ErrIdentityTaken = errors.New("identity already linked to a user")

type ExternalIdentity struct {
	Provider Provider
	Subject  string
	Email    string
	Name     string
}
