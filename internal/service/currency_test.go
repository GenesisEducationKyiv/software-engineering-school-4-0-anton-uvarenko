package service

import (
	"errors"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/stretchr/testify/assert"
)

type ConverterMock struct {
	ExpectedError error
}

func (m *ConverterMock) GetUAHToUSD() (float32, error) {
	if m.ExpectedError != nil {
		return 0, m.ExpectedError
	}

	return 13, nil
}

func TestRate(t *testing.T) {
	testTable := []struct {
		Name           string
		ExpectedError  error
		ExpectedResult float32
	}{
		{
			Name:           "OK",
			ExpectedError:  nil,
			ExpectedResult: 13,
		},
		{
			Name:           "Unexpected status code",
			ExpectedError:  pkg.ErrUnexpectedStatusCode,
			ExpectedResult: 0,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			converter := &ConverterMock{
				ExpectedError: testCase.ExpectedError,
			}
			service := NewCurrencySevice(converter)

			result, err := service.Rate()
			assert.Equal(t, result, testCase.ExpectedResult)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}
