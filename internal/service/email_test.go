package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmailService_AddEmail(t *testing.T) {
	mockRepo := new(mockEmailRepo)
	mockSender := new(mockEmailSender)
	mockConverter := new(mockCurrencyConverter)
	service := NewEmailService(mockRepo, mockSender, mockConverter)

	ctx := context.Background()

	tests := []struct {
		name          string
		email         string
		mockRepoError error
		expectedError error
	}{
		{
			name:          "Successful add email",
			email:         "test@example.com",
			mockRepoError: nil,
			expectedError: nil,
		},
		{
			name:          "Duplicate email",
			email:         "duplicate@example.com",
			mockRepoError: &pgconn.PgError{Code: "23505"},
			expectedError: pkg.ErrEmailConflict,
		},
		{
			name:          "Internal DB error",
			email:         "error@example.com",
			mockRepoError: errors.New("some db error"),
			expectedError: pkg.ErrDBInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("AddEmail", ctx, tt.email).Return(tt.mockRepoError)

			err := service.AddEmail(ctx, tt.email)

			assert.Equal(t, tt.expectedError, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEmailService_SendEmails(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name               string
		mockRepoEmails     []string
		mockRepoError      error
		mockConverterRate  float32
		mockConverterError error
		mockSenderError    error
		expectedError      error
	}{
		{
			name:               "Successful send emails",
			mockRepoEmails:     []string{"test1@example.com", "test2@example.com"},
			mockRepoError:      nil,
			mockConverterRate:  28.5,
			mockConverterError: nil,
			mockSenderError:    nil,
			expectedError:      nil,
		},
		{
			name:               "No emails registered",
			mockRepoEmails:     nil,
			mockRepoError:      sql.ErrNoRows,
			mockConverterRate:  0,
			mockConverterError: nil,
			mockSenderError:    nil,
			expectedError:      pkg.ErrNoEmailsRegistered,
		},
		{
			name:               "DB error retrieving emails",
			mockRepoEmails:     nil,
			mockRepoError:      errors.New("some db error"),
			mockConverterRate:  0,
			mockConverterError: nil,
			mockSenderError:    nil,
			expectedError:      pkg.ErrDBInternal,
		},
		{
			name:               "Currency converter error",
			mockRepoEmails:     []string{"test1@example.com"},
			mockRepoError:      nil,
			mockConverterRate:  0,
			mockConverterError: errors.New("converter error"),
			mockSenderError:    nil,
			expectedError:      errors.New("converter error"),
		},
		{
			name:               "Email send error",
			mockRepoEmails:     []string{"test1@example.com"},
			mockRepoError:      nil,
			mockConverterRate:  28.5,
			mockConverterError: nil,
			mockSenderError:    errors.New("send error"),
			expectedError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockEmailRepo)
			mockSender := new(mockEmailSender)
			mockConverter := new(mockCurrencyConverter)
			service := NewEmailService(mockRepo, mockSender, mockConverter)

			mockRepo.On("GetAll", ctx).Return(tt.mockRepoEmails, tt.mockRepoError)
			mockConverter.On("GetUAHToUSD").Return(tt.mockConverterRate, tt.mockConverterError)
			for _, email := range tt.mockRepoEmails {
				mockSender.On("SendEmail", email, mock.AnythingOfType("string")).Return(tt.mockSenderError)
			}

			err := service.SendEmails(ctx)
			t.Logf("error is: %v", err)

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.Nil(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockConverter.AssertExpectations(t)
			mockSender.AssertExpectations(t)
		})
	}
}
