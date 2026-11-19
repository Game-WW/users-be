package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/pkg/errs"
)

func (s *Server) GetAvatar(
	ctx context.Context,
	req *userspb.GetAvatarRequest,
) (*userspb.GetAvatarResponse, error) {
	if metaValue(ctx, userIDMetaKey) == "" {
		return nil, toStatus(errs.ErrUnauthenticated)
	}

	result, err := s.getAvatar.Handle(ctx, queries.GetAvatarQuery{
		UserID: req.GetUserId(),
	})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.GetAvatarResponse{Image: result.Image}, nil
}
