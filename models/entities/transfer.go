package entities

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
	CreatedAt     int64 `json:"created_at"`
}

type transferRepo interface {
	CreateTransfer(transfer *Transfer) error
	ListTransfersFromAccount(accountID int64) ([]Transfer, error)
}
