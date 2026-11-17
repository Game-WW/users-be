package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/pkg/errs"
)

func (s *Server) SetAvatarFileID(
	ctx context.Context,
	req *userspb.SetAvatarFileIDRequest,
) (*userspb.SetAvatarFileIDResponse, error) {
	if req == nil {
		return nil, toStatus(errs.ErrInvalidArgument)
	}

	err := s.setAvatarFileID.Handle(ctx, commands.SetAvatarFileIDCommand{
		UserID: metaValue(ctx, userIDMetaKey),
		FileID: req.GetFileId(),
	})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.SetAvatarFileIDResponse{}, nil
}
