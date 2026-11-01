package rpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	userIDMetaKey        = "x-user-id"
	authorizationMetaKey = "authorization"
	refreshTokenMetaKey  = "x-refresh-token"
	userAgentMetaKey     = "x-user-agent"
	forwardedForMetaKey  = "x-forwarded-for"
	xAuthCodeMetaKey     = "x-auth-code"

	accessTokenMetaKey = "x-access-token"
	logoutMetaKey      = "x-logout"
)

func metaValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	value := md.Get(key)
	if len(value) == 0 {
		return ""
	}

	return value[0]
}

func clientIP(ctx context.Context) string {
	ip, _, _ := strings.Cut(metaValue(ctx, forwardedForMetaKey), ",")
	return strings.TrimSpace(ip)
}

func sendHeader(ctx context.Context, kv ...string) error {
	return grpc.SetHeader(ctx, metadata.Pairs(kv...))
}

func bearerToken(ctx context.Context) string {
	token, ok := strings.CutPrefix(metaValue(ctx, authorizationMetaKey), "Bearer ")
	if !ok {
		return ""
	}
	
	return token
}
