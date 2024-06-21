package transport

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestSubscribe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		testName           string
		email              string
		mockAddEmailError  error
		expectedStatusCode int
	}{
		{
			testName:           "Successful subscription",
			email:              "test@example.com",
			mockAddEmailError:  nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			testName:           "Email conflict",
			email:              "conflict@example.com",
			mockAddEmailError:  pkg.ErrEmailConflict,
			expectedStatusCode: http.StatusConflict,
		},
		{
			testName:           "Internal server error",
			email:              "error@example.com",
			mockAddEmailError:  errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockEmailService := new(mockEmailService)
			mockEmailService.On("AddEmail", mock.Anything, tt.email).Return(tt.mockAddEmailError)

			handler := NewEmailHandler(mockEmailService)

			r := gin.Default()
			r.POST("/subscribe", handler.Subscribe)

			req, _ := http.NewRequest(http.MethodPost, "/subscribe", nil)
			req.PostForm = make(map[string][]string)
			req.PostForm.Add("email", tt.email)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			mockEmailService.AssertExpectations(t)
		})
	}
}
