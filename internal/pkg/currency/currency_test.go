package currency

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/stretchr/testify/assert"
)

type HTTPMock struct {
	ShouldError                bool
	ShouldUnxepectedStatusCode bool
	ShouldWrongBody            bool
}

func (m *HTTPMock) Get(url string) (resp *http.Response, err error) {
	if m.ShouldError {
		return nil, errors.New("some error")
	}

	responses := []response{
		{
			CurrencyCodeA: USDISO4217Code,
			CurrencyCodeB: UAHISO4217Code,
			RateSell:      13,
		},
	}
	result := &http.Response{
		StatusCode: http.StatusOK,
	}
	data, _ := json.Marshal(responses)
	result.Body = io.NopCloser(bytes.NewBuffer(data))

	if m.ShouldUnxepectedStatusCode {
		result.StatusCode = http.StatusInternalServerError
	}

	if m.ShouldWrongBody {
		result.Body = io.NopCloser(bytes.NewBuffer([]byte("somebody once told me")))
	}

	return result, nil
}

func TestGetUAHToUSD(t *testing.T) {
	testTable := []struct {
		Name                       string
		ExpectedError              error
		ExpectedResult             float32
		ShouldError                bool
		ShouldWrongBody            bool
		ShouldUnxepectedStatusCode bool
	}{
		{
			Name:           "OK",
			ExpectedError:  nil,
			ExpectedResult: 13,
		},
		{
			Name:                       "Unexpected status code",
			ExpectedError:              pkg.ErrUnexpectedStatusCode,
			ExpectedResult:             0,
			ShouldUnxepectedStatusCode: true,
		},
		{
			Name:            "Wrong body",
			ExpectedError:   pkg.ErrFailDecodeResponse,
			ExpectedResult:  0,
			ShouldWrongBody: true,
		},
		{
			Name:           "Perform request",
			ExpectedError:  pkg.ErrFailPerformRequest,
			ExpectedResult: 0,
			ShouldError:    true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			client := &HTTPMock{
				ShouldError:                testCase.ShouldError,
				ShouldWrongBody:            testCase.ShouldWrongBody,
				ShouldUnxepectedStatusCode: testCase.ShouldUnxepectedStatusCode,
			}
			converter := NewCurrencyConverter(client)

			result, err := converter.GetUAHToUSD()

			assert.Equal(t, testCase.ExpectedResult, result)
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}
