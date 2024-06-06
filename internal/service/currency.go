package service

import (
	"fmt"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
)

type CurrencyService struct {
	converter currencyConverter
}

type currencyConverter interface {
	GetUAHToUSD() (float32, error)
}

func NewCurrencySevice(converter currencyConverter) *CurrencyService {
	return &CurrencyService{
		converter: converter,
	}
}

func (s CurrencyService) Rate() (float32, error) {
	rate, err := s.converter.GetUAHToUSD()
	if err != nil {
		fmt.Printf("%v: [%v]", pkg.ErrRate, err)
		return 0, err
	}

	return rate, nil
}
