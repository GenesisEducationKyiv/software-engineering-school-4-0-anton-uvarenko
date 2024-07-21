package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailService struct {
	emailRepo          emailRepo
	emailEventProducer emailEventProducer
}

func NewEmailService(emailRepo emailRepo, emailEventProducer emailEventProducer) *EmailService {
	return &EmailService{
		emailRepo:          emailRepo,
		emailEventProducer: emailEventProducer,
	}
}

type emailEventProducer interface {
	ProduceEmailEvent(email string) error
}

type emailRepo interface {
	AddUser(ctx context.Context, email pgtype.Text) error
}

type rateConverter interface {
	GetUAHToUSD() (float32, error)
}

func (s *EmailService) AddEmail(ctx context.Context, email string) error {
	err := s.emailRepo.AddUser(ctx, pgtype.Text{String: email, Valid: true})
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

	err = s.emailEventProducer.ProduceEmailEvent(email)
	if err != nil {
		fmt.Printf("can't produce email event: %v", err)
	}

	return nil
}
