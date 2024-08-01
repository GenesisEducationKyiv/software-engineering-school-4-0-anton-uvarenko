package service

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"go.uber.org/zap"
)

type RateService struct {
	converter rateConverter
	logger    *zap.Logger
}

type rateConverter interface {
	GetUAHToUSD() (float32, error)
}

func NewRateservice(converter rateConverter, logger *zap.Logger) *RateService {
	return &RateService{
		converter: converter,
		logger:    logger.With(zap.String("service", "RateService")),
	}
}

func (s *RateService) GetUAHToUSD() (float32, error) {
	logger := s.logger.With(zap.String("method", "GetUAHToUSD"))

	rate, err := s.converter.GetUAHToUSD()
	if err != nil {
		logger.Error("rate error", zap.Error(err))
		return 0, fmt.Errorf("%v: %v", pkg.ErrRate, err)
	}

	return rate, nil
}
