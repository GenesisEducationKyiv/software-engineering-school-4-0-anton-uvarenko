package repo

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EmailSender struct {
	producer  *kafka.Producer
	topicName string
}

func NewEmailSender() *EmailSender {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})
	if err != nil {
		panic(err)
	}

	return &EmailSender{
		producer:  p,
		topicName: "emails",
	}
}

type emailPayload struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (s *EmailSender) SendEmail(To string, message string) error {
	payload, err := json.Marshal(emailPayload{
		To:      To,
		Message: message,
	})
	if err != nil {
		return fmt.Errorf("can't marshal payload: %w", err)
	}

	err = s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.topicName,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, nil)
	if err != nil {
		return fmt.Errorf("can't send email to kafka")
	}

	return nil
}
