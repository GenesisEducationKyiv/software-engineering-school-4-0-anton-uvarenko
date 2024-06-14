package server

import (
	"context"
	"fmt"
	"log"

	"github.com/anton-uvarenko/backend_school/internal/pkg"
	"github.com/anton-uvarenko/backend_school/internal/service"
	"github.com/go-co-op/gocron/v2"
)

func RegisterJobs(scheduler gocron.Scheduler, emailService *service.EmailService) {
	_, err := scheduler.NewJob(
		gocron.CronJob("0 3 * * *", false),
		gocron.NewTask(
			func(service *service.EmailService) {
				err := service.SendEmails(context.Background())
				if err != nil {
					fmt.Printf("%v: [%v]\n", pkg.ErrCronJob, err)
				}
			},
			emailService,
		))
	if err != nil {
		log.Fatal(err)
	}
}
