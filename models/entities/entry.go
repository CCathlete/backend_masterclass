package entities

type Entry struct {
	ID        int64
	AccountID int64
	Amount    int64
	CreatedAt int64
}

type entryRepo interface {
	CreateEntry(entry *Entry) error
	ListAccountEntries(accountID int64) ([]Entry, error)
}
