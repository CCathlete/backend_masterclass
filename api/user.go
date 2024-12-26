package api

import (
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username" binding:"required,alphanum"`
	FullName          string    `json:"full_name" binding:"required"`
	Email             string    `json:"email" binding:"required,email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user sqlc.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		Email:    req.Email,
	}
	hash, err := u.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg.HashedPassword = hash

	user, err := server.Store.CreateUser(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	res := newUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

// ------------------------------------------------------------------- //
type getUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.GetUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	ctx.JSON(http.StatusOK, user)
}

// ------------------------------------------------------------------- //
type deleteUserRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.GetUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	err = server.Store.DeleteUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	ctx.JSON(http.StatusOK, fmt.Sprintln("User ", user,
		"deleted successfully."))
}

// ------------------------------------------------------------------- //
/*
We want to display the list of users in chunks (pages). Each chunk
has a size of page_size. In order to navigate to the right place
in the whole list, we need to know how many pages to skip, this is
the offset which is the (num_of_pages_to_skip - 1) * page_size.
*/
type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listUsers(ctx *gin.Context) {
	var req listUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.Store.ListUsers(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	ctx.JSON(http.StatusOK, users)
}

// ------------------------------------------------------------------- //
type updateUserRequest struct {
	PasswordHash string `json:"password_hash" binding:"required"`
	Username     string `json:"username" binding:"required,alphanum"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.UpdateUserParams{
		HashedPassword: req.PasswordHash,
		Username:       req.Username,
	}

	userBefore, err := server.Store.GetUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	userAfter, err := server.Store.UpdateUser(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	output := struct{ Before, After sqlc.User }{
		Before: userBefore,
		After:  userAfter,
	}

	ctx.JSON(http.StatusOK, output)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	userResponse         `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {

	// Parsing the request body from the context.
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.Store.GetUser(ctx, req.Username)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
	}

	err = u.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	signedTokenString, err := server.TokenMaker.CreateToken(
		user.Username,
		server.Config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:          signedTokenString,
		AccessTokenExpiresAt: time.Now().Add(server.Config.AccessTokenDuration),
		userResponse:         newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
