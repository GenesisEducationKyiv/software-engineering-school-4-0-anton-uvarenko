package pkg

import "errors"

var (
	ErrEmailConflict = errors.New("email is already registered")
	ErrDBInternal    = errors.New("db internal error")

	ErrFailPerformRequest   = errors.New("can't perform http request")
	ErrFailDecodeResponse   = errors.New("can't decode response")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrCurrencyNotFound     = errors.New("can't find currency")

	ErrNoEmailsRegistered = errors.New("no emails are registered")
	ErrEmailSend          = errors.New("can't send email")
	ErrCronJob            = errors.New("cronjob error")

	ErrRate = errors.New("rate error")

	ErrProviders = errors.New("no successful response from providers")
)
