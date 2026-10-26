package rpc

import (
	"context"
	"strings"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
)

func (s *Server) Login(ctx context.Context, req *userspb.LoginRequest) (*userspb.LoginResponse, error) {
	result, err := s.login.Handle(ctx, commands.LoginCommand{
		Provider:  req.GetProvider(),
		IDToken:   bearerToken(ctx),
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

func bearerToken(ctx context.Context) string {
	token, ok := strings.CutPrefix(metaValue(ctx, authorizationMetaKey), "Bearer ")
	if !ok {
		return ""
	}

	return token
}
