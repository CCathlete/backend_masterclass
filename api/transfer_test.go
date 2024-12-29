package api

import (
	mockdb "backend-masterclass/db/mock"
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	u "backend-masterclass/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createReturnedResult(
	account1ID,
	account2ID,
	amount int64,
) (returnedResult sqlc.TransferTxResult) {

	returnedTransfer := sqlc.Transfer{
		ID:            1,
		FromAccountID: account1ID,
		ToAccountID:   account2ID,
		Amount:        amount,
		Currency:      u.USD,
		CreatedAt:     time.Now(),
	}

	returnedResult = sqlc.TransferTxResult{
		Transfer: returnedTransfer,
		FromEntry: sqlc.Entry{
			ID:        1,
			AccountID: account1ID,
			Amount:    -amount,
			Currency:  u.USD,
			CreatedAt: time.Now(),
		},
		ToEntry: sqlc.Entry{
			ID:        2,
			AccountID: account2ID,
			Amount:    amount,
			Currency:  u.USD,
			CreatedAt: time.Now(),
		},
	}

	return
}

// TODO: Add authorisation to every test case and call it in the loop.

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser()
	account1 := randomAccount(user1.Username)
	user2, _ := randomUser()
	account2 := randomAccount(user2.Username)
	user3, _ := randomUser()
	account3 := randomAccount(user3.Username)

	account1.Currency = u.USD
	account2.Currency = u.USD
	account3.Currency = u.ILS

	testCases := []struct {
		name      string
		setupAuth func(
			t *testing.T,
			request *http.Request,
			tokenMaker token.Maker,
			authorizationType string,
			username string,
			duration time.Duration,
		)
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
				authorizationType string,
				username string,
				duration time.Duration,
			) {
				addAuthorisation(t, request, tokenMaker, authorizationType, user1.Username, duration)
			},
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

				// Even if err = nil it'll go through TranslateError.
				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(nil, false)

				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)

				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(nil, false)

				arg := sqlc.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
					Currency:      u.USD,
				}

				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).Return(
					createReturnedResult(account1.ID, account2.ID, amount), nil)

				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(nil, false)

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransferResult(
					t,
					recorder.Body,
					createReturnedResult(account1.ID, account2.ID, amount),
				)
			},
		},
		{
			name: "from account not found",
			setupAuth: func(
				t *testing.T,
				request *http.Request,
				tokenMaker token.Maker,
				authorizationType string,
				username string,
				duration time.Duration,
			) {
				addAuthorisation(t, request, tokenMaker, authorizationType, user1.Username, duration)
			},
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// The validation function gets the from account first and checks it before getting the to account. If the from account is not found, the to account is never checked.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(sqlc.Account{}, sqlc.ErrRecordNotFound)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(sqlc.ErrRecordNotFound, true)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "to account not found",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			// The validation function gets the from account first, checks it, and only then gets and checks the to account.
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(sqlc.Account{}, sqlc.ErrRecordNotFound)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(sqlc.ErrRecordNotFound, true)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "from account currency mismatch",
			body: gin.H{
				"from_account_id": account3.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// The validation funs on the from account first, checking the currency and only if it's correct, it checks the to account.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)
				// The to account is never checked because the from account has a currency mismatch.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0).
					Return(account2, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				// --------No TranslateError since it's server level------------
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "to account currency mismatch",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// Both accounts exist but when the to account is checked, it has a currency mismatch. So the transfer is not executed.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account3.ID)).
					Times(1).
					Return(account3, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				// --------No TranslateError since it's server level------------
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "invalid currency",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        "invalid_coin",
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We won't even get to the validation function because the gin context checks the currency in the response json through a binding validator.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				// --------No TranslateError since it's server level------------
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "negative amount",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          -amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We won't even get to the validation function because the gin context checks the amount in the response json through a binding validator.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(0)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(0)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				// --------No TranslateError since it's server level------------
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccount error",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We simulate an internal error in the database. We won't get the second account because the validation function stops at the error it encounters.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(0).
					Return(sqlc.Account{}, sql.ErrConnDone)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().TranslateError(gomock.Any()).Times(1).
					Return(sqlc.ErrRecordNotFound, true)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTx error",
			body: gin.H{
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"amount":          amount,
				"currency":        u.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// We will get both accounts but the transfer will fail. The TransferTx function will return an error and the createTransfer function will write it to the response.
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
					Times(1).
					Return(account2, nil)
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sqlc.Transfer{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tokenMaker, err := token.NewPasetoMaker(u.RandomStr(32))
			require.NoError(t, err)
			config := u.Config{
				AccessTokenDuration: time.Minute,
			}

			tc.buildStubs(store)

			server := NewServer(store, config, tokenMaker)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/transfers"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data)) // We create a new request with the data we created.
			require.NoError(t, err)

			// user1 is the owner of fromAccount.
			tc.setupAuth(t, request, tokenMaker, authorisationTypeBearer, user1.Username, time.Minute)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func requireBodyMatchTransferResult(t *testing.T,
	body *bytes.Buffer,
	transferResult sqlc.TransferTxResult,
) {
	var gotResult sqlc.TransferTxResult
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)
	// --------------------Comparing transfers----------------------------
	require.Equal(
		t, gotResult.Transfer.ID, transferResult.Transfer.ID)
	require.Equal(
		t, gotResult.Transfer.Amount, transferResult.Transfer.Amount)
	require.Equal(
		t, gotResult.Transfer.FromAccountID, transferResult.Transfer.FromAccountID)
	require.Equal(
		t, gotResult.Transfer.ToAccountID, transferResult.Transfer.ToAccountID)
	require.Equal(
		t, gotResult.Transfer.Currency, transferResult.Transfer.Currency,
	)
	require.WithinDuration(t, gotResult.Transfer.CreatedAt,
		transferResult.Transfer.CreatedAt, time.Second)

	// --------------------Comparing accounts-----------------------------
	require.Equal(
		t, gotResult.FromAccount.ID, transferResult.FromAccount.ID)
	require.Equal(
		t, gotResult.FromAccount.Owner, transferResult.FromAccount.Owner)
	require.Equal(
		t, gotResult.FromAccount.Balance, transferResult.FromAccount.Balance)
	require.Equal(
		t, gotResult.FromAccount.Currency, transferResult.FromAccount.Currency)
	require.WithinDuration(t, gotResult.FromAccount.CreatedAt,
		transferResult.FromAccount.CreatedAt, time.Second)

	require.Equal(
		t, gotResult.ToAccount.ID, transferResult.ToAccount.ID)
	require.Equal(
		t, gotResult.ToAccount.Owner, transferResult.ToAccount.Owner)
	require.Equal(
		t, gotResult.ToAccount.Balance, transferResult.ToAccount.Balance)
	require.Equal(
		t, gotResult.ToAccount.Currency, transferResult.ToAccount.Currency)
	require.WithinDuration(t, gotResult.ToAccount.CreatedAt,
		transferResult.ToAccount.CreatedAt, time.Second)

	// --------------------Comparing entries-----------------------------
	require.Equal(
		t, gotResult.FromEntry.ID, transferResult.FromEntry.ID)
	require.Equal(
		t, gotResult.FromEntry.AccountID, transferResult.FromEntry.AccountID)
	require.Equal(
		t, gotResult.FromEntry.AccountID, transferResult.FromEntry.AccountID)
	require.Equal(
		t, gotResult.FromEntry.Amount, transferResult.FromEntry.Amount)
	require.Equal(
		t, gotResult.FromEntry.Currency, transferResult.FromEntry.Currency)
	require.WithinDuration(t, gotResult.FromEntry.CreatedAt,
		transferResult.FromEntry.CreatedAt, time.Second)

	require.Equal(
		t, gotResult.ToEntry.ID, transferResult.ToEntry.ID)
	require.Equal(
		t, gotResult.ToEntry.AccountID, transferResult.ToEntry.AccountID)
	require.Equal(
		t, gotResult.ToEntry.AccountID, transferResult.ToEntry.AccountID)
	require.Equal(
		t, gotResult.ToEntry.Amount, transferResult.ToEntry.Amount)
	require.Equal(
		t, gotResult.ToEntry.Currency, transferResult.ToEntry.Currency)
	require.WithinDuration(t, gotResult.ToEntry.CreatedAt,
		transferResult.ToEntry.CreatedAt, time.Second)
}
