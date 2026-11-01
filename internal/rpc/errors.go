package rpc

import (
	"errors"
	"log"

	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/pkg/errs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatus(err error) error {
	switch {
	case errors.Is(err, commands.ErrIdentityTaken):
		return status.Error(codes.AlreadyExists, "identity already linked to another user")
	case errors.Is(err, errs.ErrNotFound):
		return status.Error(codes.NotFound, "not found")
	case errors.Is(err, errs.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, "unauthenticated")
	case errors.Is(err, errs.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, "invalid argument")
	default:
		log.Printf("rpc: internal error: %v", err)
		return status.Error(codes.Internal, "internal error")
	}
}
