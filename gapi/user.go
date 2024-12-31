package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	u "backend-masterclass/util"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
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

func newUserResponse(user sqlc.User) *rpc.UserResponse {
	return &rpc.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func newCreateUserResponse(user sqlc.User) *rpc.CreateUserResponse {
	return &rpc.CreateUserResponse{
		Body: newUserResponse(user),
	}
}
