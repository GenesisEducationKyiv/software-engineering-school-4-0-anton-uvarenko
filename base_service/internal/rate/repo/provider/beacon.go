package provider

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
)

type BeaconProvider struct {
	httpClient HTTPClient
	apiKey     string
}

func NewBeaconProvider(client HTTPClient, apiKey string) *BeaconProvider {
	return &BeaconProvider{
		httpClient: client,
		apiKey:     apiKey,
	}
}

type beaconResponse struct {
	Value float32 `json:"value"`
}

func (p *BeaconProvider) GetUAHToUSD() (float32, error) {
	urlParams := url.Values{}
	urlParams.Add("api_key", p.apiKey)
	urlParams.Add("from", "USD")
	urlParams.Add("to", "UAH")
	urlParams.Add("amount", "1")

	resp, err := p.httpClient.Get("https://api.currencybeacon.com/v1?" + p.apiKey)
	if err != nil {
		return 0, pkg.ErrFailPerformRequest
	}
	defer func() {
		resp.Body.Close()
	}()

	result := beaconResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, pkg.ErrFailDecodeResponse
	}

	fmt.Printf("beacon rate: %f", result.Value)

	return result.Value, nil
}
