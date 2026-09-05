package logic

import (
	"cliwallet/internal/model"
	"errors"
	"math"
	"strconv"
)

var (
	ErrInvalidAmount          = errors.New("invalid amount")
	ErrAmountNegative         = errors.New("amount cannot be negative")
	ErrUnsupportedCurrency    = errors.New("unsupported currency")
	ErOperationWithZeroAmount = errors.New("operation cannot be performed with a 0 amount")
)

func ParseMoney(amountStr string, currency model.Currency) (model.Money, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return model.InvalidMoney, ErrInvalidAmount
	}

	if amount < 0 {
		return model.InvalidMoney, ErrAmountNegative
	}

	amountNormalized, err := NormalizeMoney(amount, currency)
	if err != nil {
		return model.InvalidMoney, err
	}

	return model.NewMoney(amountNormalized, currency), nil
}

func NormalizeMoney(amount float64, currency model.Currency) (int64, error) {
	exponent, err := currencyNormalizationExponent(currency)
	if err != nil {
		return 0, err
	}
	multiplier := math.Pow10(exponent)
	return int64(math.Round(amount * multiplier)), nil
}

func currencyNormalizationExponent(currency model.Currency) (int, error) {
	switch currency {
	case model.EUR, model.USD, model.GBP:
		return 2, nil

	case model.JPY:
		return 0, nil

	case model.KWD:
		return 3, nil

	default:
		return 0, ErrUnsupportedCurrency
	}
}
