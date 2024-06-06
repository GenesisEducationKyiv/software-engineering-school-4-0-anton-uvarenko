package service

import (
	"github.com/anton-uvarenko/backend_school/internal/pkg/currency"
	"github.com/anton-uvarenko/backend_school/internal/pkg/email"
)

type Service struct {
	EmailService    *EmailService
	CurrencyService *CurrencyService
}

func NewService(emailRepo emailRepo, emailSender *email.EmailSender, converter *currency.CurrencyConverter) *Service {
	return &Service{
		EmailService:    NewEmailService(emailRepo, emailSender, converter),
		CurrencyService: NewCurrencySevice(converter),
	}
}
