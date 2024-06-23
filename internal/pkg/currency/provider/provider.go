package provider

import "net/http"

type CurrencyProivder interface {
	GetUAHToUSD() (float32, error)
}

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}
