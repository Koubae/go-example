package logic

import (
	"testing"

	"cliwallet/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrencyNormalizationExponent(t *testing.T) {
	t.Parallel()

	t.Run("on-unsupported-currency", func(t *testing.T) {
		result, err := currencyNormalizationExponent(model.Currency("invalid"))

		require.Error(t, err)
		assert.Equal(t, 0, result)
		assert.ErrorIs(t, err, ErrUnsupportedCurrency)
	})

	tests := []struct {
		currency model.Currency
		want     int
	}{
		{model.EUR, 2},
		{model.USD, 2},
		{model.GBP, 2},
		{model.JPY, 0},
		{model.KWD, 3},
	}
	for _, tt := range tests {
		t.Run(string(tt.currency), func(t *testing.T) {

			result, err := currencyNormalizationExponent(tt.currency)

			require.NoError(t, err)
			assert.Equal(t, tt.want, result)

		})
	}

}

func TestNormalizeMoney(t *testing.T) {
	t.Parallel()

	t.Run("on-unsupported-currency", func(t *testing.T) {
		result, err := NormalizeMoney(120, model.Currency("invalid"))

		require.Error(t, err)
		assert.Equal(t, int64(0), result)
		assert.ErrorIs(t, err, ErrUnsupportedCurrency)
	})

	tests := []struct {
		name     string
		amount   float64
		currency model.Currency
		want     int64
	}{
		{
			name:     "EUR converts to cents",
			amount:   12.34,
			currency: model.EUR,
			want:     1234,
		},
		{
			name:     "USD converts to cents",
			amount:   10.50,
			currency: model.USD,
			want:     1050,
		},
		{
			name:     "GBP converts to cents",
			amount:   100,
			currency: model.GBP,
			want:     10000,
		},
		{
			name:     "JPY has no minor decimal units",
			amount:   120,
			currency: model.JPY,
			want:     120,
		},
		{
			name:     "JPY rounds fractional amount",
			amount:   120.6,
			currency: model.JPY,
			want:     121,
		},
		{
			name:     "KWD uses three decimal places",
			amount:   12.345,
			currency: model.KWD,
			want:     12345,
		},
		{
			name:     "zero amount",
			amount:   0,
			currency: model.EUR,
			want:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			result, err := NormalizeMoney(tt.amount, tt.currency)

			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestParseMoneyOnSuccess(t *testing.T) {
	t.Parallel()

	t.Run("by value", func(t *testing.T) {
		tests := []struct {
			amount   string
			currency model.Currency
			expected int64
		}{
			{"10.50", model.EUR, int64(1050)},
			{"10", model.EUR, int64(1000)},
			{"0.01", model.EUR, int64(1)},
			{"0.10", model.EUR, int64(10)},
			{"1.0", model.EUR, int64(100)},
		}
		for _, tt := range tests {
			t.Run(tt.amount, func(t *testing.T) {
				money, err := ParseMoney(tt.amount, tt.currency)

				require.NoError(t, err)
				assert.Equal(t, tt.expected, money.Amount)
				assert.True(t, money.Valid)
			})
		}
	})

	t.Run("by currency", func(t *testing.T) {
		tests := []struct {
			name     string
			amount   string
			currency model.Currency
			expected int64
		}{
			{"EUR decimal", "10.50", model.EUR, 1050},
			{"EUR integer", "10", model.EUR, 1000},
			{"EUR one cent", "0.01", model.EUR, 1},
			{"EUR zero", "0", model.EUR, 0},

			{"JPY integer", "100", model.JPY, 100},
			{"JPY rounds down", "100.4", model.JPY, 100},
			{"JPY rounds up", "100.6", model.JPY, 101},

			{"KWD three decimals", "10.123", model.KWD, 10123},
			{"KWD integer", "10", model.KWD, 10000},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				money, err := ParseMoney(tt.amount, tt.currency)

				require.NoError(t, err)
				assert.Equal(t, tt.expected, money.Amount)
				assert.Equal(t, tt.currency, money.Currency)
				assert.True(t, money.Valid)

			})
		}

	})

}

func TestParseMoneyOnError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		amount string

		currency model.Currency
		err      error
	}{
		{"amount not decimal", "10.50$", model.EUR, ErrInvalidAmount},
		{"amount is negative", "-10.50", model.EUR, ErrAmountNegative},
		{"unsupported currency", "10.50", model.Currency("invalid"), ErrUnsupportedCurrency},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			money, err := ParseMoney(tt.amount, tt.currency)

			require.Error(t, err)
			assert.ErrorIs(t, tt.err, err)
			assert.Equal(t, model.InvalidMoney, money)
		})
	}

}
