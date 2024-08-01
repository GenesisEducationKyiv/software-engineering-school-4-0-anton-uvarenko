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
		return
	}

	conn := db.Connect()
	emailDBRepo := emailRepo.New(conn)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("can't create zap logger: %v", err)
		return
	}

	kafkaProducer := producer.NewProducer(logger)
	err = kafkaProducer.RegisterTopics()
	if err != nil {
		logger.Error("can't regiser topics", zap.Error(err))
		return
	}

	emailService := emailService.NewEmailService(emailDBRepo, kafkaProducer, logger)
	emailHanlder := emailTranport.NewEmailHandler(emailService, logger)

	monobankProvider := rateRepoProvider.NewMonobankProvider(http.DefaultClient, logger)
	beaconProvider := rateRepoProvider.NewBeaconProvider(http.DefaultClient, os.Getenv("BEACONAPIKEY"), logger)
	privatProvider := rateRepoProvider.NewPrivatBankProvider(http.DefaultClient, logger)

	baseMonobankChain := rateRepoChain.NewBaseChain(monobankProvider, logger)
	baseBeaconChain := rateRepoChain.NewBaseChain(beaconProvider, logger)
	basePrivatChain := rateRepoChain.NewBaseChain(privatProvider, logger)

	baseMonobankChain.SetNext(baseBeaconChain)
	baseBeaconChain.SetNext(basePrivatChain)

	rateConverterService := rateService.NewRateservice(baseMonobankChain, logger)
	rateHandler := rateTransport.NewRateHandler(rateConverterService, logger)

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
