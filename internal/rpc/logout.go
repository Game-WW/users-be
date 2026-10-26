package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
)

func (s *Server) Logout(ctx context.Context, _ *userspb.LogoutRequest) (*userspb.LogoutResponse, error) {
	cmd := commands.LogoutCommand{
		RefreshToken: metaValue(ctx, refreshTokenMetaKey),
	}
	if err := s.logout.Handle(ctx, cmd); err != nil {
		return nil, toStatus(err)
	}

	if err := sendHeader(ctx, logoutMetaKey, "true"); err != nil {
		return nil, toStatus(err)
	}

	return &userspb.LogoutResponse{}, nil
}
