package provider

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
)

type PrivatBankProvider struct {
	httpClient HTTPClient
}

func NewPrivatBankProvider(client HTTPClient) *MonobankProvider {
	return &MonobankProvider{
		httpClient: client,
	}
}

type privatBankResponse struct {
	BaseCCY   string  `json:"base_ccy"`
	CCY       string  `json:"ccy"`
	RateSsell float32 `json:"sale"`
}

func (c *PrivatBankProvider) GetUAHToUSD() (float32, error) {
	resp, err := c.httpClient.Get("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5")
	if err != nil {
		return 0, fmt.Errorf("%v: %v", pkg.ErrFailPerformRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, pkg.ErrUnexpectedStatusCode
	}

	var result []privatBankResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("%v: %v", pkg.ErrFailDecodeResponse, err)
	}
	fmt.Printf("privatbank rates: %v", result)

	currency, err := c.findUahToUsd(result)
	if err != nil {
		return 0, err
	}

	return currency.RateSsell, nil
}

func (PrivatBankProvider) findUahToUsd(data []privatBankResponse) (privatBankResponse, error) {
	for _, v := range data {
		if v.BaseCCY == "UAH" && v.CCY == "USD" {
			return v, nil
		}
	}

	return privatBankResponse{}, pkg.ErrCurrencyNotFound
}
