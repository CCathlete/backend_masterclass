package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	"fmt"
	"time"

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

func newLoginUserResponse(
	sessionID,
	accessTokenString,
	refreshTokenString string,
	accessTokenExpiresAt,
	refreshTokenExpiresAt time.Time,
	user sqlc.User,
) *rpc.LoginUserResponse {
	return &rpc.LoginUserResponse{
		SessionId:             sessionID,
		AccessToken:           accessTokenString,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenExpiresAt),
		RefreshToken:          refreshTokenString,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenExpiresAt),
		User:                  newUserResponse(user),
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

func newDeleteUserResponse(username string) *rpc.DeleteUserResponse {
	message := fmt.Sprintf("Deleted user: %s", username)
	return &rpc.DeleteUserResponse{
		Message: message,
	}
}

func newUpdateUserResponse(before, after sqlc.User) *rpc.UpdateUserResponse {
	return &rpc.UpdateUserResponse{
		BeforeUpdate: newUserResponse(before),
		AfterUpdate:  newUserResponse(after),
	}
}
