package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
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
	ProduceSubscribedEvent(email string) error
	ProduceUnsubscribedEvent(email string) error
}

type emailRepo interface {
	AddUser(ctx context.Context, email pgtype.Text) error
	DeleteUser(ctx context.Context, email pgtype.Text) error
}

// TODO: add transactional outbox

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

	err = s.emailEventProducer.ProduceSubscribedEvent(email)
	if err != nil {
		fmt.Printf("can't produce subscribe event: %v", err)
	}

	return nil
}

func (s *EmailService) Unsubscribe(ctx context.Context, email string) error {
	err := s.emailRepo.DeleteUser(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrEmailIsNotRegistered
		}

		fmt.Printf("%v: [%v]\n", pkg.ErrDBInternal, err)
		return pkg.ErrDBInternal
	}

	err = s.emailEventProducer.ProduceUnsubscribedEvent(email)
	if err != nil {
		fmt.Printf("can't produce unsubscribe event: %v", err)
	}

	return nil
}
