package jobs

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type RateJob struct {
	scheduler    gocron.Scheduler
	emailService emailService
	logger       *zap.Logger
}

type emailService interface {
	SendEmails(ctx context.Context, rate float32) error
}

func NewRateJob(scheduler gocron.Scheduler, emailService emailService, logger *zap.Logger) *RateJob {
	return &RateJob{
		scheduler:    scheduler,
		emailService: emailService,
		logger:       logger.With(zap.String("service", "RateJob")),
	}
}

func (j *RateJob) RegisterJob() error {
	logger := j.logger.With(zap.String("method", "RegisterJob"))

	_, err := j.scheduler.NewJob(
		gocron.CronJob("", false),
		gocron.NewTask(j.startNorificationFlow),
	)
	if err != nil {
		logger.Error("can't regirster job", zap.Error(err))
	}

	return err
}

type rateResponse struct {
	Rate float32 `json:"rate"`
}

func (j *RateJob) startNorificationFlow() {
	logger := j.logger.With(zap.String("method", "startNorificationFlow"))

	resp, err := http.Get("http://localhost:8080/rate")
	if err != nil {
		logger.Error("can't perform request to rate service", zap.Error(err))
		return
	}
	defer resp.Body.Close()

	var rate rateResponse
	err = json.NewDecoder(resp.Body).Decode(&rate)
	if err != nil {
		logger.Error("ca't decode response", zap.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	err = j.emailService.SendEmails(ctx, rate.Rate)
	if err != nil {
		logger.Error("can't send email", zap.Error(err))
		return
	}
}
