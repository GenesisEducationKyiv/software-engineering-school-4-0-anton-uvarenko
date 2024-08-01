package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RateHandler struct {
	currencyService rateService
	logger          *zap.Logger
}

func NewRateHandler(currencyService rateService, logger *zap.Logger) *RateHandler {
	return &RateHandler{
		currencyService: currencyService,
		logger:          logger.With(zap.String("service", "RateHandler")),
	}
}

type rateService interface {
	GetUAHToUSD() (float32, error)
}

func (h *RateHandler) Rate(ctx *gin.Context) {
	logger := h.logger.With(zap.String("method", "Rate"))

	rate, err := h.currencyService.GetUAHToUSD()
	if err != nil {
		logger.Warn("can't retrieve rate", zap.Error(err))
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, rate)
}
