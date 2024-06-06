package transport

import (
	"net/http"

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
		// commeted because documentation doesn't expect this
		// if errors.Is(err, pkg.ErrCurrencyNotFound) {
		// 	ctx.AbortWithStatus(http.StatusNotFound)
		// 	return
		// }

		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, rate)
}
