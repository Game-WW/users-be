package queries

import (
	"context"
	"time"

	"github.com/Game-WW/users-be/internal/core/domain/user"
)

type LinkedIdentity struct {
	Provider user.Provider
	Email    *string
	LinkedAt time.Time
}

type LinkedIdentitiesReader interface {
	GetLinkedIdentities(ctx context.Context, userID string) ([]LinkedIdentity, error)
}

type ListIdentitiesHandler struct {
	reader LinkedIdentitiesReader
}

func NewListIdentitiesHandler(reader LinkedIdentitiesReader) *ListIdentitiesHandler {
	return &ListIdentitiesHandler{reader: reader}
}

func (h *ListIdentitiesHandler) Handle(ctx context.Context, q ListIdentitiesQuery) (ListIdentitiesResult, error) {
	allProviders := user.Providers()

	linkedIdentities, err := h.reader.GetLinkedIdentities(ctx, q.UserID)
	if err != nil {
		return ListIdentitiesResult{}, err
	}

	linkedIdentitiesByProvider := make(map[user.Provider]LinkedIdentity, len(linkedIdentities))
	for _, identity := range linkedIdentities {
		linkedIdentitiesByProvider[identity.Provider] = identity
	}

	result := make([]Identity, 0, len(allProviders))
	for _, provider := range allProviders {
		item := Identity{Provider: provider, Linked: false}
		if li, ok := linkedIdentitiesByProvider[provider]; ok {
			t := li.LinkedAt

			item.LinkedAt = &t
			item.Linked = true
			item.Email = li.Email
		}

		result = append(result, item)
	}

	return ListIdentitiesResult{Identities: result}, nil
}
