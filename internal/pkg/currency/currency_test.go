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
	"github.com/stretchr/testify/mock"
)

func getRateSellResponse(rateSell int) *http.Response {
	resp := []response{
		{
			CurrencyCodeA: USDISO4217Code,
			CurrencyCodeB: UAHISO4217Code,
			RateSell:      float32(rateSell),
		},
	}

	body, _ := json.Marshal(resp)

	return &http.Response{
		Body:       io.NopCloser(bytes.NewBuffer(body)),
		StatusCode: http.StatusOK,
	}
}

func TestGetUAHToUSD(t *testing.T) {
	testTable := []struct {
		Name           string
		ExpectedError  error
		ExpectedResult float32
		ClientError    error
		ClientResponse *http.Response
	}{
		{
			Name:           "OK",
			ExpectedError:  nil,
			ExpectedResult: 13,
			ClientResponse: getRateSellResponse(13),
		},
		{
			Name:           "Unexpected status code",
			ExpectedError:  pkg.ErrUnexpectedStatusCode,
			ExpectedResult: 0,
			ClientResponse: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewBuffer([]byte("something"))),
			},
		},
		{
			Name:           "Wrong body",
			ExpectedError:  pkg.ErrFailDecodeResponse,
			ExpectedResult: 0,
			ClientResponse: &http.Response{Body: io.NopCloser(bytes.NewBuffer([]byte("somebody once told me"))), StatusCode: http.StatusOK},
		},
		{
			Name:           "Perform request",
			ExpectedError:  pkg.ErrFailPerformRequest,
			ExpectedResult: 0,
			ClientError:    errors.New("some error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			client := new(MockHTTPClient)
			client.On("Get", mock.AnythingOfType("string")).Return(testCase.ClientResponse, testCase.ClientError)
			converter := NewCurrencyConverter(client)

			result, err := converter.GetUAHToUSD()

			assert.Equal(t, testCase.ExpectedResult, result)
			assert.Equal(t, testCase.ExpectedError, err)

			client.AssertExpectations(t)
		})
	}
}
