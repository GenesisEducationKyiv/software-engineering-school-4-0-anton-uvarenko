package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	currencyService rateService
}

func NewRateHandler(currencyService rateService) *RateHandler {
	return &RateHandler{
		currencyService: currencyService,
	}
}

type rateService interface {
	GetUAHToUSD() (float32, error)
}

func (h *RateHandler) Rate(ctx *gin.Context) {
	rate, err := h.currencyService.GetUAHToUSD()
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
