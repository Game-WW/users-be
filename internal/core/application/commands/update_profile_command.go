package commands

type UpdateProfileCommand struct {
	UserID           string
	Nickname         *string
	AvatarURL        *string
	Location         *string
	ShowOnlineStatus *bool
}
