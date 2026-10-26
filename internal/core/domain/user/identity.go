package user

type ExternalIdentity struct {
	Provider Provider
	Subject  string
	Email    string
	Name     string
}
