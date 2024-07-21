package producer

import (
	"context"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EmailProducer struct {
	producer  *kafka.Producer
	topicName string
}

func NewRateProducer() *EmailProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	return &EmailProducer{
		producer:  p,
		topicName: "emails",
	}
}

func (p *EmailProducer) RegisterTopics() error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	_, err = adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{
		{
			Topic:             p.topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type emailEventPayload struct {
	Email string `json:"rate"`
}

func (p *EmailProducer) ProduceEmailEvent(email string) error {
	payload, err := json.Marshal(emailEventPayload{
		Email: email,
	})
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topicName,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
