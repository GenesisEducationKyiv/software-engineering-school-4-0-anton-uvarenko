package service

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
)

type RateService struct {
	converter rateConverter
}

type rateConverter interface {
	GetUAHToUSD() (float32, error)
}

func NewRateservice(converter rateConverter) *RateService {
	return &RateService{
		converter: converter,
	}
}

func (s *RateService) GetUAHToUSD() (float32, error) {
	rate, err := s.converter.GetUAHToUSD()
	if err != nil {
		fmt.Printf("%v: [%v]", pkg.ErrRate, err)
		return 0, err
	}

	return rate, nil
}
