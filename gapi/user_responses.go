package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func newGetUserResponse(user sqlc.User) *rpc.GetUserResponse {
	return &rpc.GetUserResponse{
		Body: newUserResponse(user),
	}
}

func newListUsersResponse(users []sqlc.User) *rpc.ListUsersResponse {
	var res rpc.ListUsersResponse
	for _, user := range users {
		res.Body = append(res.Body, newUserResponse(user))
	}
	return &res
}
