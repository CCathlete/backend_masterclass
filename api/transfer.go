package api

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	// Amount greater than 0 and not min = 1 because we want to allow fractions if we'll use float insead of int in the future.
	Amount int64 `json:"amount" binding:"required,gt=0"`
	// True for both accounts (in the future we might add money conversion and allow different currencies).
	Currency string `json:"currency" binding:"required,validcurrency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}

	// -----------This includes checking user permissions.----------------
	if !server.validTransferParams(ctx, arg) {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid transfer parameters")))
		return
	}

	// After we validated the transfer parameters, we can proceed with the transfer.
	transfer, err := server.Store.TransferTx(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

// ------------------------------------------------------------------- //
type getTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTransfer(ctx *gin.Context) {
	var req getTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.Store.GetTransfer(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	// -----------Validating transfer ownership.-----------------------
	if _, ok := server.validAccountTransfer(ctx, transfer.FromAccountID, transfer.Currency, false); !ok {
		return
	}

	ctx.JSON(http.StatusOK, transfer)
}

// ------------------------------------------------------------------- //
/*
We want to display the list of transfers in chunks (pages). Each chunk
has a size of page_size. In order to navigate to the right place
in the whole list, we need to know how many pages to skip, this is
the offset which is the (num_of_pages_to_skip - 1) * page_size.
*/
type listTransfersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listTransfers(ctx *gin.Context) {
	var req listTransfersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// ---------A user can list only their own transfers.----------
	currentUser :=
		ctx.MustGet(authorisationPayloadKey).(*token.Payload).Username

	arg := sqlc.ListTransfersParams{
		// We'll list transfers from all accounts the current user owns.
		Owner:  currentUser,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.Store.ListTransfers(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

// ------------------------------------------------------------------- //
type getTransfersFromAccountRequest struct {
	AccountID int64  `json:"from_account_id" binding:"required"`
	Currency  string `json:"currency" binding:"required,validcurrency"`
}

func (server *Server) getTransfersFromAccount(ctx *gin.Context) {
	var req getTransfersFromAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// -----------------Validating account ownership.---------------------
	if _, ok :=
		server.validAccountTransfer(ctx, req.AccountID, req.Currency, false); !ok {
		return
	}

	// -----------------Getting the transfers.----------------------------
	transfers, err := server.Store.GetTransfersFrom(ctx, req.AccountID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

// ------------------------------------------------------------------- //

// ------------------------------------------------------------------- //
// TODO: Add a function to get all transfers to a specific account.
// ------------------------------------------------------------------- //

type deleteTransferRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteTransfer(ctx *gin.Context) {
	var req deleteTransferRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transfer, err := server.Store.GetTransfer(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	// TODO: This is not true, change this when we add RBAC.
	// ---------A user can delete only their own transfers.----------
	if _, ok := server.validAccountTransfer(ctx, transfer.FromAccountID, transfer.Currency, false); !ok {
		err := fmt.Errorf("user unuthorised to delete this transfer")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.Store.DeleteTransfer(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	ctx.JSON(http.StatusOK, fmt.Sprintln("Transfer ", transfer,
		"deleted successfully."))
}

// ------------------------------------------------------------------- //
type updateTransferRequest struct {
	Amount int64 `json:"amount" binding:"required"`
	ID     int64 `json:"id" binding:"required"`
}

func (server *Server) updateTransfer(ctx *gin.Context) {
	var req updateTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.UpdateTransferParams{
		Amount: req.Amount,
		ID:     req.ID,
	}

	transferBefore, err := server.Store.GetTransfer(ctx, req.ID)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	// ---------A user can update only their own transfers.----------
	if _, ok := server.validAccountTransfer(ctx, transferBefore.FromAccountID, transferBefore.Currency, false); !ok {
		err := fmt.Errorf("user unuthorised to update this transfer")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	transferAfter, err := server.Store.UpdateTransfer(ctx, arg)
	if trErr, notNil := server.Store.TranslateError(err); notNil {
		handleError(server, ctx, trErr)
		return
	}

	output := struct{ Before, After sqlc.Transfer }{
		Before: transferBefore,
		After:  transferAfter,
	}

	ctx.JSON(http.StatusOK, output)
}

// ------------------------------------------------------------------- //
