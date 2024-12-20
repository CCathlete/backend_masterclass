package u

const (
	// ILS currency.
	ILS = "ILS"
	// USD currency.
	USD = "USD"
	// EUR currency.
	EUR = "EUR"
)

// IsValidCurrency checks if the currency is valid.
func IsValidCurrency(currency string) bool {
	return currency == ILS || currency == USD || currency == EUR
}
