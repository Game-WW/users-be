package queries

import "time"

type GetMeQuery struct {
	UserID string
}

type GetMeResult struct {
	ID               string
	Nickname         *string
	AvatarURL        *string
	Location         *string
	ShowOnlineStatus bool
	CreatedAt        time.Time
}
