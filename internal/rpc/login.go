package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/domain/user"
)

func (s *Server) Login(
	ctx context.Context,
	req *userspb.LoginRequest,
	) (*userspb.LoginResponse, error) {
	provider, err := user.ParseProvider(req.Provider)
	if err != nil {
		return nil, toStatus(err)
	}

	code := metaValue(ctx, xAuthCodeMetaKey)
	codeVerifier := metaValue(ctx, xCodeVerifierMetaKey)
	authorization := bearerToken(ctx)

	creds := commands.AuthCredentials{
		Token:        authorization, // Для Google
		Code:         code,          // Для GitHub
		CodeVerifier: codeVerifier,  // Для GitHub PKCE
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
