package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyHandler_Rate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		testName           string
		mockRate           float32
		mockRateError      error
		expectedStatusCode int
		expectedBody       interface{}
	}{
		{
			testName:           "Successful rate retrieval",
			mockRate:           1.23,
			mockRateError:      nil,
			expectedStatusCode: http.StatusOK,
			expectedBody:       float32(1.23),
		},
		{
			testName:           "Currency not found",
			mockRate:           0,
			mockRateError:      pkg.ErrCurrencyNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       nil,
		},
		{
			testName:           "Internal server error",
			mockRate:           0,
			mockRateError:      errors.New("some error"),
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockCurrencyService := new(mockCurrencyService)
			mockCurrencyService.On("Rate").Return(tt.mockRate, tt.mockRateError)

			handler := NewCurrencyHandler(mockCurrencyService)

			r := gin.Default()
			r.GET("/rate", handler.Rate)

			req, _ := http.NewRequest(http.MethodGet, "/rate", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if tt.expectedStatusCode == http.StatusOK {
				var response float32
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}
			mockCurrencyService.AssertExpectations(t)
		})
	}
}
