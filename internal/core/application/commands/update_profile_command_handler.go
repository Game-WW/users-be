package commands

import (
	"context"
	"fmt"

	"github.com/Game-WW/users-be/internal/core/domain/user"
	"github.com/Game-WW/users-be/pkg/errs"
)

type ProfileUpdater interface {
	UpdateProfile(ctx context.Context, cmd UpdateProfileCommand) error
}

type UpdateProfileHandler struct {
	profiles ProfileUpdater
}

func NewUpdateProfileHandler(profiles ProfileUpdater) *UpdateProfileHandler {
	return &UpdateProfileHandler{profiles: profiles}
}

func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) error {
	if cmd.UserID == "" {
		return errs.ErrUnauthenticated
	}

	hasNoFieldsToUpdate := cmd.Nickname == nil &&
		cmd.Location == nil &&
		cmd.ShowOnlineStatus == nil
	if hasNoFieldsToUpdate {
		return nil
	}

	if err := user.Validate(cmd.Nickname, cmd.Location); err != nil {
		return err
	}

	if err := h.profiles.UpdateProfile(ctx, cmd); err != nil {
		return fmt.Errorf("update profile: %w", err)
	}

	return nil
}
