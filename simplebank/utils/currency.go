package utils

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}

// SupportedCurrencies returns the supported currencies.
func SupportedCurrencies() []string {
	return []string{USD, EUR, CAD}
}
