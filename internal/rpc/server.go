package rpc

import (
	"context"

	userspb "github.com/Game-WW/gateway-be/gen/users"
	"github.com/Game-WW/users-be/internal/core/application/commands"
	"github.com/Game-WW/users-be/internal/core/application/queries"
)

type (
	loginHandler interface {
		Handle(context.Context, commands.LoginCommand) (commands.LoginResult, error)
	}
	refreshHandler interface {
		Handle(context.Context, commands.RefreshCommand) (commands.RefreshResult, error)
	}
	logoutHandler interface {
		Handle(context.Context, commands.LogoutCommand) error
	}
	getMeHandler interface {
		Handle(context.Context, queries.GetMeQuery) (queries.GetMeResult, error)
	}
	listIdentitiesHandler interface {
		Handle(context.Context, queries.ListIdentitiesQuery) (queries.ListIdentitiesResult, error)
	}
	linkIdentityHandler interface {
		Handle(context.Context, commands.LinkIdentityCommand) (commands.LinkIdentityResult, error)
	}
	unlinkIdentityHandler interface {
		Handle(context.Context, commands.UnlinkIdentityCommand) error
	}
	updateProfileHandler interface {
		Handle(context.Context, commands.UpdateProfileCommand) error
	}
)

type Server struct {
	userspb.UnimplementedUsersServiceServer

	login          loginHandler
	refresh        refreshHandler
	logout         logoutHandler
	getMe          getMeHandler
	listIdentities listIdentitiesHandler
	linkIdentity   linkIdentityHandler
	unlinkIdentity unlinkIdentityHandler
	updateProfile  updateProfileHandler
}

var _ userspb.UsersServiceServer = (*Server)(nil)

func NewServer(
	login loginHandler,
	refresh refreshHandler,
	logout logoutHandler,
	getMe getMeHandler,
	listIdentities listIdentitiesHandler,
	linkIdentity linkIdentityHandler,
	unlinkIdentity unlinkIdentityHandler,
	updateProfile updateProfileHandler,
) *Server {
	return &Server{
		login:          login,
		refresh:        refresh,
		logout:         logout,
		getMe:          getMe,
		listIdentities: listIdentities,
		linkIdentity:   linkIdentity,
		unlinkIdentity: unlinkIdentity,
		updateProfile:  updateProfile,
	}
}
