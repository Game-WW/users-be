package commands

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/pkg/errs"
)

type AvatarFileUpdater interface {
	SetAvatarFileID(ctx context.Context, userID, fileID string) error
}

type SetAvatarFileIDHandler struct {
	profiles AvatarFileUpdater
}

func NewSetAvatarFileIDHandler(profiles AvatarFileUpdater) *SetAvatarFileIDHandler {
	return &SetAvatarFileIDHandler{profiles: profiles}
}

func (h *SetAvatarFileIDHandler) Handle(
	ctx context.Context,
	cmd SetAvatarFileIDCommand,
) error {
	if cmd.UserID == "" {
		return errs.ErrUnauthenticated
	}

	if err := h.profiles.SetAvatarFileID(ctx, cmd.UserID, cmd.FileID); err != nil {
		return fmt.Errorf("set avatar file id: %w", err)
	}

	return nil
}
