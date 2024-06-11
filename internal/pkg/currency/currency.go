package currency

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
)

type CurrencyConverter struct {
	httpClient HTTPClient
}

func NewCurrencyConverter(client HTTPClient) *CurrencyConverter {
	return &CurrencyConverter{
		httpClient: client,
	}
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type response struct {
	CurrencyCodeA int     `json:"currencyCodeA"`
	CurrencyCodeB int     `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateBuy       float32 `json:"rateBuy"`
	RateSell      float32 `json:"rateSell"`
}

const (
	UAHISO4217Code = 980
	USDISO4217Code = 840
)

func (c CurrencyConverter) GetUAHToUSD() (float32, error) {
	resp, err := c.httpClient.Get("https://api.monobank.ua/bank/currency")
	if err != nil {
		return 0, pkg.ErrFailPerformRequest
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, pkg.ErrUnexpectedStatusCode
	}

	var result []response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return 0, pkg.ErrFailDecodeResponse
	}

	currency, err := findUahToUsd(result)
	if err != nil {
		return 0, err
	}

	return currency.RateSell, nil
}

func findUahToUsd(data []response) (response, error) {
	for _, v := range data {
		if v.CurrencyCodeA == USDISO4217Code &&
			v.CurrencyCodeB == UAHISO4217Code {
			return v, nil
		}
	}

	return response{}, pkg.ErrCurrencyNotFound
}
