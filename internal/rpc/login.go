package rpc

import (
	"context"
	"fmt"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
)

func (s *Server) Login(ctx context.Context, req *userspb.LoginRequest) (*userspb.LoginResponse, error) {
	fmt.Printf("DEBUG LOGIN: provider=%s, code='%s'\n", req.GetProvider(), metaValue(ctx, "x-auth-code"))
	provider, err := user.ParseProvider(req.Provider)
	if err != nil {
		return nil, toStatus(err)
	}

	creds := commands.AuthCredentials{
		Token: metaValue(ctx, authorizationMetaKey), // Для Google
		Code:  metaValue(ctx, xAuthCodeMetaKey),   // Для GitHub
	}

	result, err := s.login.Handle(ctx, commands.LoginCommand{
		Provider:  provider,
		Creds:     creds,
		UserAgent: metaValue(ctx, userAgentMetaKey),
		IP:        clientIP(ctx),
	})
	if err != nil {
		return nil, toStatus(err)
	}

	if err := sendHeader(ctx, accessTokenMetaKey, result.AccessToken, refreshTokenMetaKey, result.RefreshToken); err != nil {
		return nil, toStatus(err)
	}

	return &userspb.LoginResponse{}, nil
}
