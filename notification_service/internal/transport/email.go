package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EmailConsumer struct {
	consumer    *kafka.Consumer
	topicName   string
	emailSender EmailSender
}

type EmailSender interface {
	SendEmail(to string, message string) error
}

func NewEmailConsumer(emailSender EmailSender) *EmailConsumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "emails",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	return &EmailConsumer{
		consumer:    consumer,
		topicName:   "emails",
		emailSender: emailSender,
	}
}

func (h EmailConsumer) InitializeTopics() {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
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
			panic(err)
		}

		h.handle(msg.Value)
	}
}

type emailPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (h EmailConsumer) handle(msg []byte) {
	payload := emailPayload{}

	err := json.Unmarshal(msg, &payload)
	if err != nil {
		fmt.Printf("can't unmarshal payload: %v", err)
		return
	}

	err = h.emailSender.SendEmail(payload.To, payload.Message)
	if err != nil {
		fmt.Printf("can't send email: %v", err)
		return
	}
}