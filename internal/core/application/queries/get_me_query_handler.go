package queries

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/pkg/errs"
)

type UserProfileReader interface {
	GetProfile(ctx context.Context, id string) (GetMeResult, error)
}

type GetMeHandler struct {
	profiles UserProfileReader
}

func NewGetMeHandler(profiles UserProfileReader) *GetMeHandler {
	return &GetMeHandler{profiles: profiles}
}

func (h *GetMeHandler) Handle(ctx context.Context, q GetMeQuery) (GetMeResult, error) {
	if q.UserID == "" {
		return GetMeResult{}, errs.ErrUnauthenticated
	}

	result, err := h.profiles.GetProfile(ctx, q.UserID)
	if err != nil {
		return GetMeResult{}, fmt.Errorf("get profile: %w", err)
	}

	return result, nil
}
