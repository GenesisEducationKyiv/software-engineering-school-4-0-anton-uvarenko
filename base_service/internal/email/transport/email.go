package transport

import (
	"context"
	"errors"
	"net/http"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/pkg"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService emailService
}

func NewEmailHandler(emailService emailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

type emailService interface {
	AddEmail(ctx context.Context, email string) error
	Unsubscribe(ctx context.Context, email string) error
}

func (h *EmailHandler) Subscribe(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")

	err := h.emailService.AddEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailConflict) {
			ctx.AbortWithStatus(http.StatusConflict)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

func (h *EmailHandler) UnSubscribe(ctx *gin.Context) {
	email := ctx.Request.FormValue("email")

	err := h.emailService.Unsubscribe(ctx, email)
	if err != nil {
		if errors.Is(err, pkg.ErrEmailIsNotRegistered) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return

	}
}
