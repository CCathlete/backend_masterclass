package api

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func isLoggedIn(ctx *gin.Context, username string) (ok bool) {
	payload, foundInCtx := ctx.Get(authorisationPayloadKey)
	if !foundInCtx {
		return false
	}

	payloadVal, ok := payload.(*token.Payload)
	if !ok {
		// ok = false.
		return
	}

	// If payload extraction was successful we can set our return value.
	ok = payloadVal.Username == username

	return
}

func (s *Server) validAccount(ctx *gin.Context, accountID int64,
	txCurrency string, isToAccount bool,
) (account sqlc.Account, ok bool) {

	account, err := s.Store.GetAccount(ctx, accountID)
	if trErr, notNil := s.Store.TranslateError(err); notNil {
		handleError(s, ctx, trErr)
	}
	if account.Currency != txCurrency {
		log.Printf("Account ID (%d) currency (%s) is different from the transfer's currency (%s).\n", accountID, account.Currency, txCurrency)
		return
	}

	// Making sure that that the logged user is the account owner or that the account is toAccount so user validation can be skipped.
	ok = isLoggedIn(ctx, account.Owner) || isToAccount
	if !ok {
		err = fmt.Errorf("account validation error (unauthorized)")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// ok = true.
	return
}

func (s *Server) validTransferParams(
	ctx *gin.Context,
	arg sqlc.TransferTxParams,
) (ok bool) {

	// ------------Checking general validity of parameters----------------
	if arg.FromAccountID == arg.ToAccountID {
		log.Println("FromAccountID and ToAccountID must be different.")
		return
	}
	if arg.Amount <= 0 {
		log.Println("Amount must be greater than 0.")
		return
	}
	if arg.FromAccountID < 1 || arg.ToAccountID < 1 {
		log.Println("Account IDs must be at least 1.")
		return
	}

	// ----------------Checking that fromAccount is valid-----------------
	_, ok = s.validAccount(ctx, arg.FromAccountID, arg.Currency, false)
	if !ok {
		log.Printf("FromAccountID (%d) is not valid.\n", arg.FromAccountID)
		// ok = false.
		return
	}

	// -------------Checking that toAccount exists & checking currencies------------------------------------------------------------
	_, ok = s.validAccount(ctx, arg.ToAccountID, arg.Currency, true)
	if !ok {
		log.Printf("ToAccountID (%d) is not valid.\n", arg.ToAccountID)
		// ok = false.
		return
	}

	return
}
