package queries

import "time"

type GetMeQuery struct {
	UserID string
}

type GetMeResult struct {
	ID          string
	UserName    *string
	Email       *string
	PhoneNumber *string
	CreatedAt   time.Time
}
