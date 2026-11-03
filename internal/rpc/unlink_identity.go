package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
)

func (s *Server) UnlinkIdentity(
	ctx context.Context,
	req *userspb.UnlinkIdentityRequest,
) (*userspb.UnlinkIdentityResponse, error) {
	provider, err := user.ParseProvider(req.GetProvider())
	if err != nil {
		return nil, toStatus(err)
	}

	err = s.unlinkIdentity.Handle(ctx, commands.UnlinkIdentityCommand{
		UserID:   metaValue(ctx, userIDMetaKey),
		Provider: provider,
	})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.UnlinkIdentityResponse{}, nil
}
