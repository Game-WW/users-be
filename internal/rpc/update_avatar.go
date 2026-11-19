package rpc

import (
	"context"
	"errors"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/avatar"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) UpdateAvatar(
	ctx context.Context,
	req *userspb.UpdateAvatarRequest,
) (*userspb.UpdateAvatarResponse, error) {
	userID := metaValue(ctx, userIDMetaKey)
	if _, err := avatar.NormalizeID(userID); err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthenticated")
	}

	err := s.uploadAvatar.Handle(ctx, commands.UploadAvatarCommand{
		UserID: userID,
		Image:  req.GetImage(),
	})

	switch {
	case err == nil:
		return &userspb.UpdateAvatarResponse{}, nil
	case errors.Is(err, avatar.ErrInvalidImage),
		errors.Is(err, avatar.ErrTooLarge),
		errors.Is(err, avatar.ErrUnsupportedFormat),
		errors.Is(err, avatar.ErrInvalidDimensions):
		return nil, status.Error(codes.InvalidArgument, err.Error())
	default:
		return nil, toStatus(err)
	}
}
