package server

import (
	"net/http"
	"time"

	"github.com/anton-uvarenko/backend_school/internal/transport"
	"github.com/gin-gonic/gin"
)

func NewServer(handler *transport.Handler) *http.Server {
	return &http.Server{
		Addr:              ":8080",
		Handler:           registerRoutes(handler),
		ReadHeaderTimeout: time.Second * 30,
	}
}

func registerRoutes(handler *transport.Handler) *gin.Engine {
	engine := gin.New()
	engine.POST("/subscribe", handler.EmailHandler.Subscribe)
	engine.GET("/rate", handler.CurrencyHandler.Rate)

	return engine
}
