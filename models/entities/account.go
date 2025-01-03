package entities

type Account struct {
	ID        int64  `json:"id"`
	Owner     string `json:"owner"`
	Balance   int64  `json:"balance"`
	Currency  string `json:"currency"`
	CreatedAt string `json:"created_at"`
}

type accountRepo interface {
	CreateAccount(account Account) (Account, error)
	GetAccount(owner string) (Account, error)
	ListAccounts() ([]Account, error)
	UpdateAccount(account Account) (Account, error)
	DeleteAccount(owner string) error
}
