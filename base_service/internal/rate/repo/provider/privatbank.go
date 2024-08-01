package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"go.uber.org/zap"
)

type PrivatBankProvider struct {
	httpClient HTTPClient
	logger     *zap.Logger
}

func NewPrivatBankProvider(client HTTPClient, logger *zap.Logger) *MonobankProvider {
	return &MonobankProvider{
		httpClient: client,
		logger:     logger.With(zap.String("service", "PrivatBankProvider")),
	}
}

type privatBankResponse struct {
	BaseCCY   string  `json:"base_ccy"`
	CCY       string  `json:"ccy"`
	RateSsell float32 `json:"sale"`
}

func (c *PrivatBankProvider) GetUAHToUSD() (float32, error) {
	logger := c.logger.With(zap.String("method", "GetUAHToUSD"))

	resp, err := c.httpClient.Get("https://api.privatbank.ua/p24api/pubinfo?exchange&coursid=5")
	if err != nil {
		logger.Error("can't perform request to api", zap.Error(err))
		return 0, fmt.Errorf("%v: %v", pkg.ErrFailPerformRequest, err)
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

	var result []privatBankResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("can't decode response", zap.Error(err))
		return 0, fmt.Errorf("%v: %v", pkg.ErrFailDecodeResponse, err)
	}

	logger.Info("rates", zap.Any("data", result))
	currency, err := c.findUahToUsd(result)
	if err != nil {
		logger.Error("no expected currency found in response", zap.Error(err))

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
