package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/jobs"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo/sender"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/service"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/transport"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/transport/consumer"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	connection := db.Connect()
	emailRepo := repo.New(connection)

	emailSender := sender.NewEmailSender(os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"))

	emailService := service.NewEmailService(emailRepo, emailSender)

	emailHandler := transport.NewEmailHandler(emailService)

	kafkaConsumer := consumer.NewConsumer(emailHandler)
	kafkaConsumer.InitializeTopics()
	go kafkaConsumer.StartPolling()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	rateJob := jobs.NewRateJob(scheduler, emailService)
	rateJob.RegisterJob()
	scheduler.Start()

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish

	err = scheduler.Shutdown()
	if err != nil {
		fmt.Printf("can't properly shutdown job scheduler: %v\n", err)
	}

	err = connection.Close(context.Background())
	if err != nil {
		fmt.Printf("can't properly close db connection: %v\n", err)
	}
}
