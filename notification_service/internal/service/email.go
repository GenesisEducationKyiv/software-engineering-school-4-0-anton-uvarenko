package service

import (
	"context"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
)

type EmailService struct {
	emailRepo   emailRepo
	emailSender emailSender
}

type emailRepo interface {
	AddEmail(ctx context.Context, arg repo.AddEmailParams) error
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

func (s *EmailService) SendEmail(to string, message string) error {
	err := s.emailSender.SendEmail(to, message)
	if err != nil {
		return err
	}

	return nil
}
