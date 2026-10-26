package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
)

func (s *Server) Refresh(ctx context.Context, _ *userspb.RefreshRequest) (*userspb.RefreshResponse, error) {
	cmd := commands.RefreshCommand{
		RefreshToken: metaValue(ctx, refreshTokenMetaKey),
	}
	result, err := s.refresh.Handle(ctx, cmd)
	if err != nil {
		return nil, toStatus(err)
	}

	if err := sendHeader(ctx, accessTokenMetaKey, result.AccessToken); err != nil {
		return nil, toStatus(err)
	}

	return &userspb.RefreshResponse{}, nil
}
