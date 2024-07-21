package service

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
)

type RateService struct {
	converter         rateConverter
	rateEventProducer rateEventProducer
}

type rateConverter interface {
	GetUAHToUSD() (float32, error)
}

type rateEventProducer interface {
	ProduceRateEvent(rate float32) error
}

func NewRateSevice(converter rateConverter, rateEventProducer rateEventProducer) *RateService {
	return &RateService{
		converter:         converter,
		rateEventProducer: rateEventProducer,
	}
}

func (s *RateService) GetUAHToUSD() (float32, error) {
	rate, err := s.converter.GetUAHToUSD()
	if err != nil {
		fmt.Printf("%v: [%v]", pkg.ErrRate, err)
		return 0, err
	}

	err = s.rateEventProducer.ProduceRateEvent(rate)
	if err != nil {
		fmt.Printf("can't produce rate event: %v", err)
	}

	return rate, nil
}
