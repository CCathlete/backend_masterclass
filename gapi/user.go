package gapi

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/rpc"
	u "backend-masterclass/util"
	"context"
	"log"
)

// ------------------------------------------------------------------- //
// ------------------------------------------------------------------- //

func (server *Server) CreateUser(
	ctx context.Context,
	req *rpc.CreateUserRequest,
) (res *rpc.CreateUserResponse, err error) {

	// -------------------Validating the request.-------------------------
	violations := validateCreateUserRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

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
// ------------------------------------------------------------------- //

func (server *Server) LoginUser(
	ctx context.Context,
	req *rpc.LoginUserRequest,
) (res *rpc.LoginUserResponse, err error) {

	username := req.GetUsername()

	// -------------------Validating the request.-------------------------
	violations := validateLoginUserRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

	// -------------------Getting user's details.-------------------------
	user, err := server.Store.GetUser(ctx, username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// ---------------Verifying the password.-----------------------------
	err = u.CheckPassword(req.GetPassword(), user.HashedPassword)
	if err != nil {
		err = handleError(err)
		return
	}

	// ---------------Deleting existing session.--------------------------
	// ---------------Will take affect only if a session exists.----------
	err = server.Store.DeleteSessionByUsername(ctx, username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}

	// -----------Creating a signed authentication token.-----------------
	accessTokenString, accessPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.Config.AccessTokenDuration,
	)
	if err != nil {
		err = handleError(err)
		return
	}

	// ----------------Creating a refresh token.--------------------------
	refreshTokenString, refreshPayload, err := server.TokenMaker.CreateToken(
		user.Username,
		server.Config.RefreshTokenDuration,
	)
	if err != nil {
		err = handleError(err)
		return
	}

	// ----------------Saving the refresh token.--------------------------
	metadata := server.extractMetadata(ctx)
	session, err := server.Store.CreateSession(ctx, sqlc.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshTokenString,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		err = handleError(trErr)
		return
	}
	log.Printf("Created session: %+v\n", session)

	// ----------------Setting up the response.---------------------------
	res = newLoginUserResponse(
		session.ID.String(),
		accessTokenString,
		refreshTokenString,
		accessPayload.ExpiresAt,
		refreshPayload.ExpiresAt,
		user,
	)
	return
}

// ------------------------------------------------------------------- //
// ------------------------------------------------------------------- //
func (server *Server) GetUser(
	ctx context.Context,
	req *rpc.GetUserRequest,
) (res *rpc.GetUserResponse, err error) {

	// -------------------Validating the request.-------------------------
	violations := validateGetUserRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

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
// ------------------------------------------------------------------- //

func (server *Server) ListUsers(
	ctx context.Context,
	req *rpc.ListUsersRequest,
) (res *rpc.ListUsersResponse, err error) {

	// -------------------Validating the request.-------------------------
	violations := validateListUsersRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

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
// ------------------------------------------------------------------- //

func (server *Server) DeleteUser(
	ctx context.Context,
	req *rpc.DeleteUserRequest,
) (res *rpc.DeleteUserResponse, err error) {

	// -------------------Validating the request.-------------------------
	violations := validateDeleteUserRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

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

// ------------------------------------------------------------------- //
// ------------------------------------------------------------------- //

func (server *Server) UpdateUser(
	ctx context.Context,
	req *rpc.UpdateUserRequest,
) (res *rpc.UpdateUserResponse, err error) {

	// -------------------Validating the request.-------------------------
	violations := validateUpdateUserRequest(req)
	if violationsFound(violations) {
		err = badRequestError(violations)
		return
	}

	// ----------------Setting up parameters for the query.---------------
	arg := sqlc.UpdateUserParams{
		HashedPassword: req.GetPasswordHash(),
		Username:       req.GetUsername(),
	}

	// -------------------Executing the query.----------------------------
	userBefore, err := server.Store.GetUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(trErr)
		return
	}

	userAfter, err := server.Store.UpdateUser(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(trErr)
		return
	}

	// ----------------Setting up the response.---------------------------
	res = newUpdateUserResponse(userBefore, userAfter)
	return
}
