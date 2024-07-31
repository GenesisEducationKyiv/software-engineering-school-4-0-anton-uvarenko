package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo/sender"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type EmailService struct {
	emailRepo   emailRepo
	emailSender emailSender
	logger      *zap.Logger
}

type emailRepo interface {
	AddEmail(ctx context.Context, arg repo.AddEmailParams) error
	GetAll(ctx context.Context) ([]repo.SendedEmail, error)
	DeleteEmail(ctx context.Context, email pgtype.Text) error
	UpdateEmail(ctx context.Context, arg repo.UpdateEmailParams) error
}

type emailSender interface {
	SendEmail(to string, message string) error
}

func NewEmailService(emailRepo emailRepo, emailSender emailSender, logger *zap.Logger) *EmailService {
	return &EmailService{
		emailRepo:   emailRepo,
		emailSender: emailSender,
		logger:      logger.With(zap.String("service", "EmailService")),
	}
}

func (s *EmailService) SaveEmail(ctx context.Context, arg repo.AddEmailParams) error {
	logger := s.logger.With(zap.String("method", "SaveEmail"))

	err := s.emailRepo.AddEmail(ctx, arg)
	if err != nil {
		logger.Error("can't add email to db", zap.Error(err))
		return err
	}

	return nil
}

func (s *EmailService) DeleteEmail(ctx context.Context, email string) error {
	logger := s.logger.With(zap.String("method", "DeleteEmail"))

	err := s.emailRepo.DeleteEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		logger.Error("can't delete email from db", zap.Error(err))
		return err
	}
	return nil
}

func (s *EmailService) SendEmails(ctx context.Context, rate float32) error {
	logger := s.logger.With(zap.String("method", "SendEmails"))

	sendedEmails, err := s.emailRepo.GetAll(ctx)
	if err != nil {
		logger.Error("can't retrieve emails", zap.Error(err))
		return err
	}

	if len(sendedEmails) == 0 {
		logger.Warn("sened emails is empty")
		return errors.New("no sended emails")
	}

	wg := &sync.WaitGroup{}

	for _, sendedEmail := range sendedEmails {
		wg.Add(1)

		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			err := s.SendEmail(sendedEmail.Email.String, fmt.Sprintf("%s %f", sender.DefaultEmailMessage, rate))
			if err != nil {
				// log email cause don't have id
				logger.Error("can't send email", zap.String("id", sendedEmail.Email.String), zap.Error(err))
				return
			}

			err = s.emailRepo.UpdateEmail(ctx, repo.UpdateEmailParams{
				Email:     pgtype.Text{String: sendedEmail.Email.String, Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
			})
			if err != nil {
				logger.Error("can't update sended email", zap.String("id", sendedEmail.Email.String), zap.Error(err))
				return
			}
		}(wg)
	}

	wg.Wait()

	return nil
}

func (s *EmailService) SendEmail(to string, message string) error {
	logger := s.logger.With(zap.String("method", "SendEmail"))

	err := s.emailSender.SendEmail(to, message)
	if err != nil {
		logger.Error("can't send email", zap.Error(err))
		return err
	}

	return nil
}
