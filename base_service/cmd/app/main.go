package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/db"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/server"
	"github.com/joho/godotenv"

	emailCore "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/email/core"
	emailRepo "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/email/repo"
	emailService "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/email/service"
	emailTranport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/email/transport"

	rateRepoChain "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/rate/repo/chain"
	rateRepoProvider "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/rate/repo/provider"
	rateService "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/rate/service"
	rateTransport "github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/base_sevice/internal/rate/transport"
)

func main() {
	godotenv.Load()

	conn := db.Connect()
	emailDBRepo := emailCore.New(conn)

	monobankProvider := rateRepoProvider.NewMonobankProvider(http.DefaultClient)
	beaconProvider := rateRepoProvider.NewBeaconProvider(http.DefaultClient, os.Getenv("BEACONAPIKEY"))
	privatProvider := rateRepoProvider.NewPrivatBankProvider(http.DefaultClient)

	baseMonobankChain := rateRepoChain.NewBaseChain(monobankProvider)
	baseBeaconChain := rateRepoChain.NewBaseChain(beaconProvider)
	basePrivatChain := rateRepoChain.NewBaseChain(privatProvider)

	baseMonobankChain.SetNext(baseBeaconChain)
	baseBeaconChain.SetNext(basePrivatChain)

	rateConverterService := rateService.NewRateSevice(baseMonobankChain)

	emailSender := emailRepo.NewEmailSender()
	emailService := emailService.NewEmailService(emailDBRepo, emailSender, rateConverterService)
	emailHanlder := emailTranport.NewEmailHandler(emailService)
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
