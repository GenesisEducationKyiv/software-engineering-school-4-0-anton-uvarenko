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
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/transport/consumer"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	connection := db.Connect()
	emailRepo := repo.New(connection)

	emailSender := sender.NewEmailSender(os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"))

	emailService := service.NewEmailService(emailRepo, emailSender)

	emailHandler := transport.NewEmailHandler(emailService)

	rateHandler := transport.NewRateHandler(emailService)

	kafkaConsumer := consumer.NewConsumer(rateHandler, emailHandler)
	kafkaConsumer.InitializeTopics()
	go kafkaConsumer.StartPolling()

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish
}
