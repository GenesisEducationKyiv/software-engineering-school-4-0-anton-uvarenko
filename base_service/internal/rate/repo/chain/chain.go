package chain

import (
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/repo/provider"
	"go.uber.org/zap"
)

type ProvidersChain interface {
	GetUAHToUSD() (float32, error)
	SetNext(next ProvidersChain)
}

type BaseChain struct {
	provider provider.CurrencyProvider
	next     ProvidersChain
	logger   *zap.Logger
}

func NewBaseChain(provider provider.CurrencyProvider, logger *zap.Logger) *BaseChain {
	return &BaseChain{
		provider: provider,
		logger:   logger.With(zap.String("service", "BaseChain")),
	}
}

func (c *BaseChain) SetNext(next ProvidersChain) {
	c.next = next
}

func (c *BaseChain) GetUAHToUSD() (float32, error) {
	rate, err := c.provider.GetUAHToUSD()
	if err == nil {
		c.logger.Error("error receiving rate", zap.Error(err))
		return rate, nil
	}

	next := c.next
	if next == nil {
		return 0, pkg.ErrProviders
	}

	return next.GetUAHToUSD()
}
