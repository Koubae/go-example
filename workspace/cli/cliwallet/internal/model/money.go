package model

type Currency string

type Money struct {
	Amount   int64
	Currency Currency
	Valid    bool
}

func NewMoney(amount int64, currency Currency) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
		Valid:    true,
	}
}

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	KWD Currency = "KWD"
)

var (
	InvalidMoney = Money{}
)
