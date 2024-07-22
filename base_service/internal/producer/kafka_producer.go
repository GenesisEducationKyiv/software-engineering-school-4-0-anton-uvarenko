package producer

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topics   []string
}

func NewProducer() *Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		panic(err)
	}

	return &Producer{
		producer: p,
		topics:   []string{"emails", "rate"},
	}
}

func (p *Producer) RegisterTopics() error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		fmt.Println("can't create admin client")
		panic(err)
	}

	topicSpecifications := []kafka.TopicSpecification{}
	for _, topicName := range p.topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{Topic: topicName, NumPartitions: 1, ReplicationFactor: 1})
	}

	_, err = adminClient.CreateTopics(context.Background(), topicSpecifications)
	if err != nil {
		fmt.Println("can't register topics")

		return err
	}

	return nil
}
