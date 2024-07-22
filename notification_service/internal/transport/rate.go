package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RateHandler struct {
	emailService rateConsumerEmailService
}

type rateConsumerEmailService interface {
	SendEmails(ctx context.Context, rate float32) error
}

func NewRateHandler(emailService rateConsumerEmailService) *RateHandler {
	return &RateHandler{
		emailService: emailService,
	}
}

type ratePayload struct {
	Rate float32 `json:"rate"`
}

func (h RateHandler) Handle(msg []byte) error {
	fmt.Println("starting rate handling")
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)

	payload := ratePayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return err
	}

	fmt.Println("starting sending emails")
	err = h.emailService.SendEmails(ctx, payload.Rate)
	if err != nil {
		fmt.Printf("can't send email: %v", err)
		return err
	}

	return nil
}
