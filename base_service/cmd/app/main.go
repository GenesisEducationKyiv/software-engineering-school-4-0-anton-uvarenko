package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/producer"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/server"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	emailRepo "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/email/repo"
	emailService "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/email/service"
	emailTranport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/email/transport"

	rateRepoChain "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/repo/chain"
	rateRepoProvider "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/repo/provider"
	rateService "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/service"
	rateTransport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_service/internal/rate/transport"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("can't load env: %v", err)
		panic(err)
	}

	conn := db.Connect()
	emailDBRepo := emailRepo.New(conn)

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	kafkaProducer := producer.NewProducer(logger)
	err = kafkaProducer.RegisterTopics()
	if err != nil {
		fmt.Printf("can't register topics: %v", err)
		panic(err)
	}

	emailService := emailService.NewEmailService(emailDBRepo, kafkaProducer, logger)
	emailHanlder := emailTranport.NewEmailHandler(emailService, logger)

	monobankProvider := rateRepoProvider.NewMonobankProvider(http.DefaultClient)
	beaconProvider := rateRepoProvider.NewBeaconProvider(http.DefaultClient, os.Getenv("BEACONAPIKEY"))
	privatProvider := rateRepoProvider.NewPrivatBankProvider(http.DefaultClient)

	baseMonobankChain := rateRepoChain.NewBaseChain(monobankProvider)
	baseBeaconChain := rateRepoChain.NewBaseChain(beaconProvider)
	basePrivatChain := rateRepoChain.NewBaseChain(privatProvider)

	baseMonobankChain.SetNext(baseBeaconChain)
	baseBeaconChain.SetNext(basePrivatChain)

	rateConverterService := rateService.NewRateservice(baseMonobankChain)
	rateHandler := rateTransport.NewRateHandler(rateConverterService)

	httpServer := server.NewServer(rateHandler, emailHanlder)
	go log.Fatal(httpServer.ListenAndServe())

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, syscall.SIGTERM)

	<-finish
	func() {
		err := conn.Close(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}()
}
