package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo/sender"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailService struct {
	emailRepo   emailRepo
	emailSender emailSender
}

type emailRepo interface {
	AddEmail(ctx context.Context, arg repo.AddEmailParams) error
	GetAll(ctx context.Context) ([]repo.SendedEmail, error)
	UpdateEmail(ctx context.Context, arg repo.UpdateEmailParams) error
}

type emailSender interface {
	SendEmail(to string, message string) error
}

func NewEmailService(emailRepo emailRepo, emailSender emailSender) *EmailService {
	return &EmailService{
		emailRepo:   emailRepo,
		emailSender: emailSender,
	}
}

func (s *EmailService) SaveEmail(ctx context.Context, arg repo.AddEmailParams) error {
	err := s.emailRepo.AddEmail(ctx, arg)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailService) SendEmails(ctx context.Context, rate float32) error {
	sendedEmails, err := s.emailRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(sendedEmails) == 0 {
		return nil
	}

	sendedCount := &atomic.Int32{}
	wg := &sync.WaitGroup{}

	for _, sendedEmail := range sendedEmails {
		wg.Add(1)

		go func() {
			if time.Since(sendedEmail.UpdatedAt.Time) < time.Hour*24 {
				return
			}

			err := s.SendEmail(sendedEmail.Email.String, fmt.Sprintf("%s %f", sender.DEFAULT_EMAIL_MESSAGE, rate))
			if err != nil {
				fmt.Printf("can't send email to %s: %v", sendedEmail.Email.String, err)
				return
			}

			err = s.emailRepo.UpdateEmail(ctx, repo.UpdateEmailParams{
				Email:     pgtype.Text{String: sendedEmail.Email.String, Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			})
			if err != nil {
				fmt.Printf("can't update email %s: %v", sendedEmail.Email.String, err)
				return
			}

			sendedCount.Add(1)
		}()
	}

	wg.Wait()

	if sendedCount.Load() == 0 {
		return errors.New("no emails were sended")
	}

	return nil
}

func (s *EmailService) SendEmail(to string, message string) error {
	err := s.emailSender.SendEmail(to, message)
	if err != nil {
		return err
	}

	return nil
}
