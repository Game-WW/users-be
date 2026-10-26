package commands

type RefreshCommand struct {
	RefreshToken string
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}
