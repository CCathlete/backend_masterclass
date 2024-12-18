package api_test

import (
	"backend-masterclass/db/sqlc"
	u "backend-masterclass/util"
	"testing"
)

func TestGetAccountAPI(t *testing.T) {

}

func randomAccount() sqlc.Account {
	return sqlc.Account{
		ID:       u.RandomInt(1, 100),
		Owner:    u.RandomOwner(),
		Balance:  u.RandomMoney(),
		Currency: u.RandCurrency(),
	}
}
