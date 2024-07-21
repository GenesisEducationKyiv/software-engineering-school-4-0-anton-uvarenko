package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type RateProducer struct {
	producer  *kafka.Producer
	topicName string
}

func NewRateProducer() *RateProducer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		panic(err)
	}

	return &RateProducer{
		producer:  p,
		topicName: "rate",
	}
}

func (p *RateProducer) RegisterTopics() error {
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

type rateEventPayload struct {
	Rate      float32   `json:"rate"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *RateProducer) ProduceRateEvent(rate float32) error {
	payload, err := json.Marshal(rateEventPayload{
		Rate:      rate,
		UpdatedAt: time.Now(),
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
