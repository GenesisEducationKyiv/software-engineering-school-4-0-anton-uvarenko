package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailHandler struct {
	emailService emailService
}

type emailService interface {
	SaveEmail(ctx context.Context, arg repo.AddEmailParams) error
	DeleteEmail(ctx context.Context, email string) error
}

func NewEmailHandler(emaemailService emailService) *EmailHandler {
	return &EmailHandler{
		emailService: emaemailService,
	}
}

type emailSubscribePayload struct {
	Email string `json:"email"`
}

func (h EmailHandler) Handle(msg *kafka.Message) error {
	originHeaderIndex := slices.IndexFunc(msg.Headers, func(header kafka.Header) bool {
		return header.Key == "origin"
	})

	if originHeaderIndex < 0 {
		return errors.New("no origin header")
	}

	var err error
	switch string(msg.Headers[originHeaderIndex].Value) {
	case "user_unsubscribed":
		h.handleUnsubscribeEvent(msg.Value)
		break
	case "user_subscribed":
		h.handleSubscribeEvent(msg.Value)
		break
	default:
		err = errors.New("unsuported origin")
	}

	return err
}

func (h EmailHandler) handleSubscribeEvent(msg []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)

	payload := emailSubscribePayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return err
	}

	err = h.emailService.SaveEmail(ctx, repo.AddEmailParams{
		Email: pgtype.Text{
			String: payload.Email,
			Valid:  true,
		},
	})
	if err != nil {
		fmt.Printf("can't save email: %v", err)
		return err
	}

	return nil
}

type emailUnsubscribePayload struct {
	Email string `json:"email"`
}

func (h EmailHandler) handleUnsubscribeEvent(msg []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)

	payload := emailUnsubscribePayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return err
	}

	err = h.emailService.DeleteEmail(ctx, payload.Email)
	if err != nil {
		fmt.Printf("can't delete email: %v", err)
		return err
	}

	return nil
}
