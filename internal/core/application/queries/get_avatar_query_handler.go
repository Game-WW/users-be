package queries

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/Game-WW/users-be/internal/core/domain/avatar"
	"github.com/Game-WW/users-be/internal/infrastructure/s3"
	"github.com/Game-WW/users-be/pkg/errs"
)

type GetAvatarHandler struct {
	storage *s3.Storage
}

func NewGetAvatarHandler(storage *s3.Storage) *GetAvatarHandler {
	return &GetAvatarHandler{storage: storage}
}

func (h *GetAvatarHandler) Handle(
	ctx context.Context,
	q GetAvatarQuery,
) (GetAvatarResult, error) {
	id, err := avatar.NormalizeID(q.UserID)
	if err != nil {
		return GetAvatarResult{}, errs.ErrInvalidArgument
	}

	key := avatar.Key(id)
	obj, _, err := h.storage.Get(ctx, key)
	if errors.Is(err, s3.ErrObjectNotFound) {
		return GetAvatarResult{}, nil
	}
	if err != nil {
		return GetAvatarResult{}, fmt.Errorf("get avatar: %w", err)
	}
	defer obj.Close()

	image, err := io.ReadAll(obj)
	if err != nil {
		return GetAvatarResult{}, fmt.Errorf("read avatar: %w", err)
	}

	return GetAvatarResult{Image: image}, nil
}
