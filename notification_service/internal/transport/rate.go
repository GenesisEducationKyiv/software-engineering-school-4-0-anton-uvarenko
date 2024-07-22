package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type RateConsumer struct {
	consumer     *kafka.Consumer
	topicName    string
	emailService rateConsumerEmailService
}

type rateConsumerEmailService interface {
	SendEmails(ctx context.Context, rate float32) error
}

func NewRateConsuer(emailService rateConsumerEmailService) *RateConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
		"group.id":          "emails",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	return &RateConsumer{
		consumer:     consumer,
		topicName:    "rate",
		emailService: emailService,
	}
}

func (h RateConsumer) InitializeTopics() {
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

func (h RateConsumer) Consume() {
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

type ratePayload struct {
	Rate float32 `json:"rate"`
}

func (h RateConsumer) handle(msg []byte) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)

	payload := ratePayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return
	}

	err = h.emailService.SendEmails(ctx, payload.Rate)
	if err != nil {
		fmt.Printf("can't send email: %v", err)
		return
	}
}
