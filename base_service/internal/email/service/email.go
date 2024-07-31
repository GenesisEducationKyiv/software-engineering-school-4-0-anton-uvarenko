package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

type EmailService struct {
	emailRepo          emailRepo
	emailEventProducer emailEventProducer
	logger             *zap.Logger
}

func NewEmailService(emailRepo emailRepo, emailEventProducer emailEventProducer, logger *zap.Logger) *EmailService {
	return &EmailService{
		emailRepo:          emailRepo,
		emailEventProducer: emailEventProducer,
		logger:             logger.With(zap.String("service", "EmailService")),
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
	logger := s.logger.With(zap.String("method", "AddEmail"))

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

		logger.Error("can't add user to db", zap.Error(err))
		return pkg.ErrDBInternal
	}

	err = s.emailEventProducer.ProduceSubscribedEvent(email)
	if err != nil {
		logger.Warn("can't produce subscribe event", zap.Error(err))
	}

	return nil
}

func (s *EmailService) Unsubscribe(ctx context.Context, email string) error {
	logger := s.logger.With(zap.String("method", "Unsubscribe"))

	err := s.emailRepo.DeleteUser(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return pkg.ErrEmailIsNotRegistered
		}

		logger.Error("can't delete user", zap.Error(err))
		return pkg.ErrDBInternal
	}

	err = s.emailEventProducer.ProduceUnsubscribedEvent(email)
	if err != nil {
		logger.Warn("can't produce unsubscribe event", zap.Error(err))
	}

	return nil
}
