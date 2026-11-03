package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetMe(
	ctx context.Context,
	_ *userspb.GetMeRequest,
	) (*userspb.GetMeResponse, error) {
	result, err := s.getMe.Handle(ctx, queries.GetMeQuery{UserID: metaValue(ctx, userIDMetaKey)})
	if err != nil {
		return nil, toStatus(err)
	}

	return &userspb.GetMeResponse{
		User: toProtoUser(result),
	}, nil
}

func toProtoUser(r queries.GetMeResult) *userspb.User {
	return &userspb.User{
		Id:          r.ID,
		Username:    r.UserName,
		Email:       r.Email,
		PhoneNumber: r.PhoneNumber,
		CreatedAt:   timestamppb.New(r.CreatedAt),
	}
}
