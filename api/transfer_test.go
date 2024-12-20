package api

import (
	mockdb "backend-masterclass/db/mock"
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createRandomAccount(t *testing.T) sqlc.Account {
	arg := sqlc.CreateAccountParams{
		Owner:    u.RandomOwner(),
		Balance:  u.RandomMoney(),
		Currency: u.RandCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	account1.Currency = u.USD
	account2.Currency = u.USD
	account3.Currency = u.ILS

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				arg := sqlc.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
					Currency:      u.USD,
				}
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransfer(t, recorder.Body, sqlc.Transfer{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
					Currency:      u.USD,
				})
			},
		},
	}

}

func requireBodyMatchTransfer(t *testing.T,
	body *bytes.Buffer,
	transfer sqlc.Transfer,
) {
	var gotTransfer sqlc.Transfer
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	err = json.Unmarshal(data, &gotTransfer)
	require.NoError(t, err)
	require.Equal(t, transfer, gotTransfer)
}
