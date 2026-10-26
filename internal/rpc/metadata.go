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

	accessTokenMetaKey = "x-access-token"
	logoutMetaKey      = "x-logout"
)

// metaValue возвращает первое значение ключа или "", если его нет.
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

// clientIP берёт первый адрес из x-forwarded-for — это клиент, остальные — прокси по пути.
func clientIP(ctx context.Context) string {
	ip, _, _ := strings.Cut(metaValue(ctx, forwardedForMetaKey), ",")
	return strings.TrimSpace(ip)
}

// sendHeader отправляет metadata ответа — gateway превращает её в cookie.
func sendHeader(ctx context.Context, kv ...string) error {
	return grpc.SetHeader(ctx, metadata.Pairs(kv...))
}
