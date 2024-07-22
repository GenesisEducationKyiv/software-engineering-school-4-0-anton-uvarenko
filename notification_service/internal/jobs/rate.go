package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type RateJob struct {
	scheduler    gocron.Scheduler
	emailService emailService
}

type emailService interface {
	SendEmails(ctx context.Context, rate float32) error
}

func NewRateJob(scheduler gocron.Scheduler, emailService emailService) *RateJob {
	return &RateJob{
		scheduler:    scheduler,
		emailService: emailService,
	}
}

func (j *RateJob) RegisterJob() {
	j.scheduler.NewJob(
		gocron.CronJob("", false),
		gocron.NewTask(j.startNorificationFlow),
	)
}

type rateResponse struct {
	Rate float32 `json:"rate"`
}

func (j *RateJob) startNorificationFlow() {
	resp, err := http.Get("http://localhost:8080/rate")
	if err != nil {
		fmt.Printf("can't perform request: %v", err)
		return
	}
	defer resp.Body.Close()

	var rate rateResponse
	err = json.NewDecoder(resp.Body).Decode(&rate)
	if err != nil {
		fmt.Printf("can't decode response: %v", err)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)

	err = j.emailService.SendEmails(ctx, rate.Rate)
	if err != nil {
		fmt.Printf("can't send emails: %v", err)
		return
	}
}
