package server

import (
	"net/http"
	"time"

	emailTransport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/email/transport"
	rateTransport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/transport"
	"github.com/gin-gonic/gin"
)

func NewServer(
	rateHandler *rateTransport.RateHandler,
	emailHandler *emailTransport.EmailHandler,
) *http.Server {
	return &http.Server{
		Addr:              ":8080",
		Handler:           registerRoutes(rateHandler, emailHandler),
		ReadHeaderTimeout: time.Second * 30,
	}
}

func registerRoutes(
	rateHandler *rateTransport.RateHandler,
	emailHandler *emailTransport.EmailHandler,
) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.POST("/subscribe", emailHandler.Subscribe)
	engine.GET("/rate", rateHandler.Rate)

	return engine
}
