package commands

type LoginCommand struct {
	Provider  string
	IDToken   string
	UserAgent string
	IP        string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}
