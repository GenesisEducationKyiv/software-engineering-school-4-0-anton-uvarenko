package provider

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"go.uber.org/zap"
)

type MonobankProvider struct {
	httpClient HTTPClient
	logger     *zap.Logger
}

func NewMonobankProvider(client HTTPClient, logger *zap.Logger) *MonobankProvider {
	return &MonobankProvider{
		httpClient: client,
		logger:     logger.With(zap.String("service", "MonobankProvider")),
	}
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

func (c *MonobankProvider) GetUAHToUSD() (float32, error) {
	logger := c.logger.With(zap.String("method", "GetUAHToUSD"))

	resp, err := c.httpClient.Get("https://api.monobank.ua/bank/currency")
	if err != nil {
		logger.Error("can't perform request to api", zap.Error(err))
		return 0, pkg.ErrFailPerformRequest
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("can't read body after bad status code", zap.Error(err))
		}

		logger.Warn(
			"unexpected status code from api",
			zap.Int("status code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return 0, pkg.ErrUnexpectedStatusCode
	}

	var result []response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("can't decode response", zap.Error(err))
		return 0, pkg.ErrFailDecodeResponse
	}

	logger.Info("rates", zap.Any("data", result))

	currency, err := c.findUahToUsd(result)
	if err != nil {
		logger.Error("no expected currency found in response", zap.Error(err))
		return 0, err
	}

	return currency.RateSell, nil
}

func (MonobankProvider) findUahToUsd(data []response) (response, error) {
	for _, v := range data {
		if v.CurrencyCodeA == USDISO4217Code &&
			v.CurrencyCodeB == UAHISO4217Code {
			return v, nil
		}
	}

	return response{}, pkg.ErrCurrencyNotFound
}
