package transport

import "github.com/anton-uvarenko/backend_school/internal/service"

type Handler struct {
	EmailHandler    *EmailHandler
	CurrencyHandler *CurrencyHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		EmailHandler:    NewEmailHandler(service.EmailService),
		CurrencyHandler: NewCurrencyHandler(service.CurrencyService),
	}
}
