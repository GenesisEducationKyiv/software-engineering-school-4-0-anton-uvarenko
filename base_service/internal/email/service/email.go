package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
	"github.com/jackc/pgx/v5/pgconn"
)

type EmailService struct {
	emailRepo         emailRepo
	emailSender       emailSender
	currencyConverter rateConverter
}

func NewEmailService(emailRepo emailRepo, sender emailSender, converter rateConverter) *EmailService {
	return &EmailService{
		emailSender:       sender,
		emailRepo:         emailRepo,
		currencyConverter: converter,
	}
}

type emailSender interface {
	SendEmail(to string, message string) error
}

type emailRepo interface {
	AddEmail(ctx context.Context, email string) error
	GetAll(ctx context.Context) ([]string, error)
}

type rateConverter interface {
	GetUAHToUSD() (float32, error)
}

func (s *EmailService) AddEmail(ctx context.Context, email string) error {
	err := s.emailRepo.AddEmail(ctx, email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// duplicate key
			err := err.(*pgconn.PgError)
			if err.Code == "23505" {
				return pkg.ErrEmailConflict
			}
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDBInternal, err)
		return pkg.ErrDBInternal
	}

	return nil
}

func (s *EmailService) SendEmails(ctx context.Context) error {
	emails, err := s.emailRepo.GetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrNoEmailsRegistered
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDBInternal, err)
		return pkg.ErrDBInternal
	}

	currentCurrency, err := s.currencyConverter.GetUAHToUSD()
	if err != nil {
		return err
	}
	for _, email := range emails {
		err = s.emailSender.SendEmail(email, fmt.Sprintf("current ratio uah to usd is %v", currentCurrency))
		if err != nil {
			fmt.Printf("%v: [%v]", pkg.ErrEmailSend, err)
		}
	}

	return nil
}
