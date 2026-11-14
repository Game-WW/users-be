package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
)

func (s *Server) LinkIdentity(
	ctx context.Context,
	req *userspb.LinkIdentityRequest,
) (*userspb.LinkIdentityResponse, error) {
	provider, err := user.ParseProvider(req.GetProvider())
	if err != nil {
		return nil, toStatus(err)
	}

	userID := metaValue(ctx, userIDMetaKey)

	creds := commands.AuthCredentials{
		Token:        bearerToken(ctx),
		Code:         metaValue(ctx, xAuthCodeMetaKey),
		CodeVerifier: metaValue(ctx, xCodeVerifierMetaKey),
	}

	result, err := s.linkIdentity.Handle(ctx, commands.LinkIdentityCommand{
		UserID:   userID,
		Provider: provider,
		Creds:    creds,
	})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.LinkIdentityResponse{
		Identity: &userspb.Identity{
			Provider: string(result.Provider),
			Linked:   true,
			Email:    &result.Email,
		},
	}, nil
}
