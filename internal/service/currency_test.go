package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyService_Rate(t *testing.T) {
	mockConverter := new(mockCurrencyConverter)
	service := NewCurrencySevice(mockConverter)

	tests := []struct {
		name          string
		mockRate      float32
		mockError     error
		expectedRate  float32
		expectedError error
	}{
		{
			name:          "Successful rate retrieval",
			mockRate:      28.5,
			mockError:     nil,
			expectedRate:  28.5,
			expectedError: nil,
		},
		{
			name:          "Currency converter error",
			mockRate:      0,
			mockError:     errors.New("converter error"),
			expectedRate:  0,
			expectedError: errors.New("converter error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConverter.On("GetUAHToUSD").Return(tt.mockRate, tt.mockError)

			rate, err := service.Rate()

			assert.Equal(t, tt.expectedRate, rate)
			assert.Equal(t, tt.expectedError, err)
			mockConverter.AssertExpectations(t)
		})
	}
}
