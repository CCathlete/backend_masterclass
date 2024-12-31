package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	u "backend-masterclass/util"
	"context"
)

func (server *Server) CreateUser(
	ctx context.Context,
	req *rpc.CreateUserRequest,
) (res *rpc.CreateUserResponse, err error) {

	// ----------------Setting up parameters for the query.---------------
	arg := sqlc.CreateUserParams{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
	}
	hash, err := u.HashPassword(req.GetPassword())
	if err != nil {
		err = handleError(err)
		return
	}
	arg.HashedPassword = hash

	// -------------------Executing the query.----------------------------
	user, err := server.Store.CreateUser(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// ----------------Setting up the response.---------------------------
	res = newCreateUserResponse(user)
	return
}

// ------------------------------------------------------------------- //

func (server *Server) GetUser(
	ctx context.Context,
	req *rpc.GetUserRequest,
) (res *rpc.GetUserResponse, err error) {

	// -------------------Executing the query.----------------------------
	user, err := server.Store.GetUser(ctx, req.GetUsername())
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// ----------------Setting up the response.---------------------------
	res = newGetUserResponse(user)
	return
}

// ------------------------------------------------------------------- //

func (server *Server) ListUsers(
	ctx context.Context,
	req *rpc.ListUsersRequest,
) (res *rpc.ListUsersResponse, err error) {

	// ----------------Setting up parameters for the query.---------------
	arg := sqlc.ListUsersParams{
		Limit:  req.GetPageSize(),
		Offset: (req.GetPageId() - 1) * req.GetPageSize(),
	}

	// -------------------Executing the query.----------------------------
	users, err := server.Store.ListUsers(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// ----------------Setting up the response.---------------------------
	res = newListUsersResponse(users)
	return
}

// ------------------------------------------------------------------- //

func (server *Server) DeleteUser(
	ctx context.Context,
	req *rpc.DeleteUserRequest,
) (res *rpc.DeleteUserResponse, err error) {

	// -------------------Executing the query.----------------------------
	err = server.Store.DeleteUser(ctx, req.GetUsername())
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// ----------------Setting up the response.---------------------------
	res = newDeleteUserResponse(req.GetUsername())
	return
}
