package commands

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/avatar"
)

type AvatarStorage interface {
	Put(ctx context.Context, key string, data []byte, contentType string) error
	Delete(ctx context.Context, key string) error
}

type UploadAvatarHandler struct {
	storage AvatarStorage
}

func NewUploadAvatarHandler(storage AvatarStorage) *UploadAvatarHandler {
	return &UploadAvatarHandler{
		storage: storage,
	}
}

func (h *UploadAvatarHandler) Handle(ctx context.Context, cmd UploadAvatarCommand) error {
	userID, err := avatar.NormalizeID(cmd.UserID)
	if err != nil {
		return err
	}

	key := avatar.Key(userID)

	if len(cmd.Image) == 0 {
		if err := h.storage.Delete(ctx, key); err != nil {
			return fmt.Errorf("delete avatar: %w", err)
		}
		return nil
	}

	contentType, err := avatar.Validate(cmd.Image)
	if err != nil {
		return err
	}

	if err := h.storage.Put(ctx, key, cmd.Image, contentType); err != nil {
		return fmt.Errorf("save avatar: %w", err)
	}

	return nil
}
