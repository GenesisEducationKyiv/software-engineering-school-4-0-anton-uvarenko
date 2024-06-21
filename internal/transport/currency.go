package transport

import (
	"errors"
	"net/http"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/gin-gonic/gin"
)

type CurrencyHandler struct {
	currencyService currencyService
}

func NewCurrencyHandler(currencyService currencyService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
	}
}

type currencyService interface {
	Rate() (float32, error)
}

func (h *CurrencyHandler) Rate(ctx *gin.Context) {
	rate, err := h.currencyService.Rate()
	if err != nil {
		if errors.Is(err, pkg.ErrCurrencyNotFound) {
			ctx.AbortWithStatus(http.StatusNotFound)
			return
		}

		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, rate)
}
