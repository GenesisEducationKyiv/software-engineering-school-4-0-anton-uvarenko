package email

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	gomail "gopkg.in/mail.v2"
)

func TestEmailSender_SendEmail(t *testing.T) {
	tests := []struct {
		name          string
		to            string
		message       string
		mockError     error
		expectedError error
	}{
		{
			name:          "Successful send",
			to:            "recipient@example.com",
			message:       "Test message",
			mockError:     nil,
			expectedError: nil,
		},
		{
			name:          "Failed send",
			to:            "recipient@example.com",
			message:       "Test message",
			mockError:     assert.AnError,
			expectedError: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDialer := new(mockDialer)
			sender := &EmailSender{
				from:     "test@example.com",
				smtpHost: "smtp.gmail.com",
				smtpPort: "587",
				password: "password",
				dialer:   mockDialer,
			}

			m := gomail.NewMessage()
			m.SetHeader("From", sender.from)
			m.SetHeader("To", tt.to)
			m.SetBody("text/plain", tt.message)
			mockDialer.On("DialAndSend", mock.Anything).Return(tt.mockError)

			err := sender.SendEmail(tt.to, tt.message)
			assert.Equal(t, tt.expectedError, err)
			mockDialer.AssertExpectations(t)
		})
	}
}
