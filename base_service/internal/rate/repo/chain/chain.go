package chain

import (
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/rate/repo/provider"
)

type ProvidersChain interface {
	GetUAHToUSD() (float32, error)
	SetNext(next ProvidersChain)
}

type BaseChain struct {
	provider provider.CurrencyProvider
	next     ProvidersChain
}

func NewBaseChain(provider provider.CurrencyProvider) *BaseChain {
	return &BaseChain{
		provider: provider,
	}
}

func (c *BaseChain) SetNext(next ProvidersChain) {
	c.next = next
}

func (c *BaseChain) GetUAHToUSD() (float32, error) {
	rate, err := c.provider.GetUAHToUSD()
	if err == nil {
		fmt.Printf("error receiving rate: %v\n", err)
		return rate, nil
	}

	next := c.next
	if next == nil {
		return 0, pkg.ErrProviders
	}

	return next.GetUAHToUSD()
}
