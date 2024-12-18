package api_test

import (
	"backend-masterclass/api"
	mockdb "backend-masterclass/db/mock"
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		accountId     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				// Checking response.
				require.Equal(t, http.StatusOK, recorder.Code)

				// Check response's body.
				requireBodyMatchAccount(t, account, recorder.Body)
			},
		},
		//TODO: Add more cases.
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)
	// build stubs.

	// Starting test server and sending request.
	server := api.NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/accounts/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	// We use the inner ServeHTTP method and not Gin's Run method because
	// we want to use our recorder as a response writer.
	server.Router.ServeHTTP(recorder, request)
}

func randomAccount() sqlc.Account {
	return sqlc.Account{
		ID:       u.RandomInt(1, 100),
		Owner:    u.RandomOwner(),
		Balance:  u.RandomMoney(),
		Currency: u.RandCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T,
	account sqlc.Account,
	body *bytes.Buffer,
) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount sqlc.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}
