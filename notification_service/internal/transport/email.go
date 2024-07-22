package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-anton-uvarenko/notification_service/internal/repo"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgtype"
)

type EmailConsumer struct {
	consumer     *kafka.Consumer
	topicName    string
	emailService emailService
}

type emailService interface {
	SaveEmail(ctx context.Context, arg repo.AddEmailParams) error
}

func NewEmailConsumer(emaemailService emailService) *EmailConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
		"group.id":          "emails",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	return &EmailConsumer{
		consumer:     consumer,
		topicName:    "emails",
		emailService: emaemailService,
	}
}

func (h EmailConsumer) InitializeTopics() {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		panic(err)
	}

	_, err = adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{
		{
			Topic:             h.topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	})
	if err != nil {
		panic(err)
	}
}

func (h EmailConsumer) Consume() {
	h.consumer.SubscribeTopics([]string{h.topicName}, nil)

	for {
		msg, err := h.consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("can't read message: %v", err)
			continue
		}

		h.handle(msg.Value)
	}
}

type emailSubscribePayload struct {
	Email string `json:"email"`
}

func (h EmailConsumer) handle(msg []byte) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)

	payload := emailSubscribePayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return
	}

	err = h.emailService.SaveEmail(ctx, repo.AddEmailParams{
		Email: pgtype.Text{
			String: payload.Email,
			Valid:  true,
		},
	})
	if err != nil {
		fmt.Printf("can't send email: %v", err)
		return
	}
}
