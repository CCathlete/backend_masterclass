package api

import (
	"backend-masterclass/db/sqlc"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	// Amount greater than 0 and not min = 1 because we want to allow fractions if we'll use float insead of int in the future.
	Amount   int64 `json:"amount" binding:"required,gt=0"`
	Currency int64 `json:"currency" binding:"required,oneof=ILS USD EUR"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	transfer, err := server.store.GetTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	arg := sqlc.ListTransfersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	transfers, err := server.store.ListTransfers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, transfers)
}

// ------------------------------------------------------------------- //
// TODO: Add a function to get all transfers from a specific account.
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

	transfer, err := server.store.GetTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}

	err = server.store.DeleteTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	transferBefore, err := server.store.GetTransfer(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
	}

	transferAfter, err := server.store.UpdateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	output := struct{ Before, After sqlc.Transfer }{
		Before: transferBefore,
		After:  transferAfter,
	}

	ctx.JSON(http.StatusOK, output)
}
