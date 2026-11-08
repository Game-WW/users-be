package commands

type UpdateProfileCommand struct {
	UserID           string
	Nickname         *string
	Location         *string
	ShowOnlineStatus *bool
}
