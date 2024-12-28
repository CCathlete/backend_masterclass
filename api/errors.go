package api

import (
	"backend-masterclass/db/sqlc"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(server *Server, ctx *gin.Context, err error) {
	if trErr, notNil := server.Store.TranslateError(err); notNil {

		if errors.Is(err, sqlc.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnprocessableEntity, errorResponse(err))
			return

		} else if errors.Is(trErr, sqlc.ErrConnection) {
			ctx.JSON(http.StatusServiceUnavailable, errorResponse(err))
			return
		}

		// Any other error.
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
}

func errorResponse(err error) (resBody gin.H) {
	return gin.H{
		"error": err.Error(),
	}
}
