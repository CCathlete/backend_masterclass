package api

import (
	"backend-masterclass/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// --------gin.Context keys (are stored in our ctx): --------------
	authorisationHeaderKey  = "authorization"
	authorisationPayloadKey = "authorization_payload"

	// --------Authorisation types (lower case): -----------------------------------
	authorisationTypeBearer = "bearer"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorisation header: auth type prefix (Bearer for now) + the token we get after logging in.
		authorisationHeader := ctx.GetHeader(authorisationHeaderKey)
		if len(authorisationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, errorResponse(err))
			return
		}

		// We want to extract the signed token string from the header.
		fields := strings.Fields(authorisationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, errorResponse(err))
			return
		}

		// We want to extract the authorisation type from the auth header.
		// Authorisation type = token type.
		authorisationType := strings.ToLower(fields[0])
		if authorisationType != authorisationTypeBearer {
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, errorResponse(err))
			return
		}

		signedTokenString := fields[1]

		payload, err := tokenMaker.VerifyToken(signedTokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorisationPayloadKey, payload)

		ctx.Next()
	}
}
