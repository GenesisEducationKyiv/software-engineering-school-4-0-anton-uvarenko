package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo/sender"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/transport"
)

func main() {
	connection := db.Connect()
	emailRepo := repo.New(connection)

	emailSender := sender.NewEmailSender(os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"))

	emailService := service.NewEmailService(emailRepo, emailSender)

	emailConsumer := transport.NewEmailConsumer(emailService)
	emailConsumer.InitializeTopics()
	go emailConsumer.Consume()

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish
}
