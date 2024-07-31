package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"go.uber.org/zap"
)

type BeaconProvider struct {
	httpClient HTTPClient
	apiKey     string
	logger     *zap.Logger
}

func NewBeaconProvider(client HTTPClient, apiKey string, logger *zap.Logger) *BeaconProvider {
	return &BeaconProvider{
		httpClient: client,
		apiKey:     apiKey,
		logger:     logger.With(zap.String("service", "BeaconProvider")),
	}
}

type beaconResponse struct {
	Value float32 `json:"value"`
}

func (p *BeaconProvider) GetUAHToUSD() (float32, error) {
	logger := p.logger.With(zap.String("method", "GetUAHToUSD"))

	urlParams := url.Values{}
	urlParams.Add("api_key", p.apiKey)
	urlParams.Add("from", "USD")
	urlParams.Add("to", "UAH")
	urlParams.Add("amount", "1")

	resp, err := p.httpClient.Get("https://api.currencybeacon.com/v1?" + p.apiKey)
	if err != nil {
		logger.Error("can't perform request to api", zap.Error(err))
		return 0, pkg.ErrFailPerformRequest
	}
	defer func() {
		resp.Body.Close()
	}()

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

	result := beaconResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Error("can't decode response", zap.Error(err))
		return 0, pkg.ErrFailDecodeResponse
	}

	logger.Info("rates", zap.Any("data", result))

	return result.Value, nil
}
