package api

import (
	"backend-masterclass/db/sqlc"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(ctx *gin.Context, err error) {

	if errors.Is(err, sqlc.ErrRecordNotFound) {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
		return

	} else if errors.Is(err, sqlc.ErrConnection) {
		ctx.JSON(http.StatusServiceUnavailable, errorResponse(err))
		return

	} else if errors.Is(err, sqlc.ErrForbiddenInput) {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// Any other error.
	ctx.JSON(http.StatusInternalServerError, errorResponse(err))
}

func errorResponse(err error) (resBody gin.H) {
	return gin.H{
		"error": err.Error(),
	}
}
