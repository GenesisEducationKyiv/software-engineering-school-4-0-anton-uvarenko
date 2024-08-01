package transport

import (
	"context"
	"errors"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmailHandler struct {
	emailService emailService
	logger       *zap.Logger
}

func NewEmailHandler(emailService emailService, logger *zap.Logger) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
		logger:       logger.With(zap.String("service", "EmailHandler")),
	}
}

type emailService interface {
	AddEmail(ctx context.Context, email string) error
	Unsubscribe(ctx context.Context, email string) error
}

func (h *EmailHandler) Subscribe(ctx *gin.Context) {
	logger := h.logger.With(zap.String("method", "Subscribe"))

	email := ctx.Request.FormValue("email")

	err := h.emailService.AddEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailConflict) {
			// i'm logging email, cause i don't have user id
			logger.Warn("email already exists", zap.String("id", email))
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error("can't subscribe user", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *EmailHandler) UnSubscribe(ctx *gin.Context) {
	logger := h.logger.With(zap.String("method", "Subscribe"))

	email := ctx.Request.FormValue("email")

	err := h.emailService.Unsubscribe(ctx, email)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailIsNotRegistered) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		logger.Error("can't unsubscribe user", zap.Error(err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return

	}
}
