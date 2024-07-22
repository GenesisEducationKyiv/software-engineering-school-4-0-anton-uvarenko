package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailHandler struct {
	emailService emailService
}

type emailService interface {
	SaveEmail(ctx context.Context, arg repo.AddEmailParams) error
}

func NewEmailHandler(emaemailService emailService) *EmailHandler {
	return &EmailHandler{
		emailService: emaemailService,
	}
}

type emailSubscribePayload struct {
	Email string `json:"email"`
}

func (h EmailHandler) Handle(msg []byte) error {
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
		fmt.Printf("can't send email: %v", err)
		return err
	}

	return nil
}
