package service

import (
	"context"
	"errors"
	"testing"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/stretchr/testify/assert"
)

type EmailSenderMock struct {
	ExpectedError error
}

func (m *EmailSenderMock) SendEmail(To string, message string) error {
	return m.ExpectedError
}

type EmailRepoMock struct{}

func (m *EmailRepoMock) AddEmail(ctx context.Context, email string) error {
	return nil
}

func (m *EmailRepoMock) GetAll(ctx context.Context) ([]string, error) {
	return []string{"a", "b", "c"}, nil
}

func TestAddEmail(t *testing.T) {
	testTable := []struct {
		Name          string
		ExpectedError error
	}{
		{
			Name:          "OK",
			ExpectedError: nil,
		},
		{
			Name:          "Some error",
			ExpectedError: pkg.ErrUnexpectedStatusCode,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.Name, func(t *testing.T) {
			sender := &EmailSenderMock{
				ExpectedError: testCase.ExpectedError,
			}

			converter := &ConverterMock{
				ExpectedError: testCase.ExpectedError,
			}

			repo := &EmailRepoMock{}

			service := NewEmailService(repo, sender, converter)

			err := service.SendEmails(context.Background())
			assert.Equal(t, errors.Is(err, testCase.ExpectedError), true)
		})
	}
}
