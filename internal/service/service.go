package service

import (
	"github.com/anton-uvarenko/backend_school/internal/pkg/currency/provider"
	"github.com/anton-uvarenko/backend_school/internal/pkg/email"
)

type Service struct {
	EmailService    *EmailService
	CurrencyService *CurrencyService
}

func NewService(emailRepo emailRepo, emailSender *email.EmailSender, converter provider.CurrencyProvider) *Service {
	return &Service{
		EmailService:    NewEmailService(emailRepo, emailSender, converter),
		CurrencyService: NewCurrencySevice(converter),
	}
}
