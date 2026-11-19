package queries

type GetAvatarQuery struct {
	UserID string
}

type GetAvatarResult struct {
	Image []byte
}
