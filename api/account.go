package api

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,validcurrency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// The middleware extracts this from the token in the auth header
	// and puts it in the context.
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)
	arg := sqlc.CreateAccountParams{
		// The owner is the username from the authorisation payload.
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.Store.CreateAccount(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {

		if errors.Is(trErr, sqlc.ErrForbiddenInput) {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return

		} else if errors.Is(trErr, sqlc.ErrConnection) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// Any other error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// ------------------------------------------------------------------- //
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Store.GetAccount(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {

		if errors.Is(err, sqlc.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return

		} else if errors.Is(trErr, sqlc.ErrConnection) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		// Any other error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Making sure that the logged in user is allowed to see the account.
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := fmt.Errorf("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// ------------------------------------------------------------------- //
// Selects account + locks the selected rows until the transaction is
// committed (e.g. suited for transaction use).
type getAccountForUpdateRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountForUpdate(ctx *gin.Context) {
	var req getAccountForUpdateRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.Store.GetAccountForUpdate(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	// Making sure that the logged in user is allowed to see the account.
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := fmt.Errorf("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

// ------------------------------------------------------------------- //
/*
We want to display the list of accounts in chunks (pages). Each chunk
has a size of page_size. In order to navigate to the right place
in the whole list, we need to know how many pages to skip, this is
the offset which is the (num_of_pages_to_skip - 1) * page_size.
*/
type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Making sure that the logged in user is allowed to see the account.
	authPayload := ctx.MustGet(authorisationPayloadKey).(*token.Payload)

	arg := sqlc.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.Store.ListAccounts(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

// ------------------------------------------------------------------- //
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Making sure that the logged in user is allowed to delete the account.
	account, ok := server.validAccount(ctx, req.ID)
	if !ok {
		return
	}

	err := server.Store.DeleteAccount(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintln("Account ", account,
		"deleted successfully."))
}

// ------------------------------------------------------------------- //
type updateAccountBalanceRequest struct {
	Amount int64 `json:"amount" binding:"required"`
	ID     int64 `json:"id" binding:"required"`
}

func (server *Server) updateAccountBalance(ctx *gin.Context) {
	var req updateAccountBalanceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Making sure that the logged in user is allowed to update the account.
	accountBefore, ok := server.validAccount(ctx, req.ID)
	if !ok {
		// Response handling is done inside the validation if not valid.
		return
	}

	arg := sqlc.UpdateAccountBalanceParams{
		Amount: req.Amount,
		ID:     req.ID,
	}

	accountAfter, err := server.Store.UpdateAccountBalance(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	output := struct{ Before, After sqlc.Account }{
		Before: accountBefore,
		After:  accountAfter,
	}

	ctx.JSON(http.StatusOK, output)
}

// ------------------------------------------------------------------- //
// type updateAccountRequest struct {
// 	Balance int64 `json:"balance" binding:"required"`
// 	ID      int64 `json:"id" binding:"required"`
// }

// func (server *Server) updateAccount(ctx *gin.Context) {
// 	var req updateAccountRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// Making sure that the logged in user is allowed to update the account.
// 	_ = ctx.MustGet(authorisationPayloadKey).(*token.Payload)

// 	arg := sqlc.UpdateAccountParams{
// 		Balance: req.Balance,
// 		ID:      req.ID,
// 	}

// 	accountBefore, err := server.Store.GetAccount(ctx, req.ID)
// 	if trErr, notNil := server.Store.TranslateError(err); notNil {
// 		handleError(server, ctx, trErr)
// return
// 	}

// 	accountAfter, err := server.Store.UpdateAccount(ctx, arg)
// 	if trErr, notNil := server.Store.TranslateError(err); notNil {
// 		handleError(server, ctx, trErr)
// return
// 	}

// 	output := struct{ Before, After sqlc.Account }{
// 		Before: accountBefore,
// 		After:  accountAfter,
// 	}

// 	ctx.JSON(http.StatusOK, output)
// }
