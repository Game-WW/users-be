package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"github.com/Game-WW/users-be/pkg/errs"
)

func (s *Server) UpdateProfile(ctx context.Context, req *userspb.UpdateProfileRequest) (*userspb.UpdateProfileResponse, error) {
	if req == nil {
		return nil, toStatus(errs.ErrInvalidArgument)
	}

	userID := metaValue(ctx, userIDMetaKey)
	err := s.updateProfile.Handle(ctx, commands.UpdateProfileCommand{
		UserID:           userID,
		Nickname:         req.Nickname,
		AvatarURL:        req.AvatarUrl,
		Location:         req.Location,
		ShowOnlineStatus: req.ShowOnlineStatus,
	})
	if err != nil {
		return nil, toStatus(err)
	}

	profile, err := s.getMe.Handle(ctx, queries.GetMeQuery{UserID: userID})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.UpdateProfileResponse{
		User: toProtoUser(profile),
	}, nil
}
