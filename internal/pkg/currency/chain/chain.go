package chain

import "github.com/anton-uvarenko/backend_school/internal/pkg/currency/provider"

type ProvidersChain interface {
	GetUAHToUSD() (float32, error)
	SetNext(next ProvidersChain)
}

type BaseChain struct {
	provider provider.CurrencyProivder
	next     ProvidersChain
}

func NewBaseChain(provider provider.CurrencyProivder) *BaseChain {
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
		return rate, nil
	}

	next := c.next
	if next == nil {
		return 0, nil
	}

	return next.GetUAHToUSD()
}
