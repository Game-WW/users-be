package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/queries"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) ListIdentities(
	ctx context.Context,
	_ *userspb.ListIdentitiesRequest,
) (*userspb.ListIdentitiesResponse, error) {
	result, err := s.listIdentities.Handle(ctx, queries.ListIdentitiesQuery{
		UserID: metaValue(ctx, userIDMetaKey),
	})
	if err != nil {
		return nil, toStatus(err)
	}

	identities := make([]*userspb.Identity, 0, len(result.Identities))
	for _, item := range result.Identities {
		identities = append(identities, toProtoIdentity(item))
	}

	return &userspb.ListIdentitiesResponse{Identities: identities}, nil
}

func toProtoIdentity(item queries.Identity) *userspb.Identity {
	p := &userspb.Identity{
		Provider: string(item.Provider),
		Linked:   item.Linked,
		Email:    item.Email,
	}
	if item.LinkedAt != nil {
		p.LinkedAt = timestamppb.New(*item.LinkedAt)
	}

	return p
}
