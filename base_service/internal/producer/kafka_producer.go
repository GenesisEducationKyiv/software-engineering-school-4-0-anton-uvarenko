package producer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

type Producer struct {
	producer *kafka.Producer
	topics   []string
	logger   *zap.Logger
}

func NewProducer(logger *zap.Logger) *Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		panic(err)
	}

	return &Producer{
		producer: p,
		topics:   []string{"emails"},
		logger:   logger.With(zap.String("service", "Producer")),
	}
}

func (p *Producer) RegisterTopics() error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		p.logger.Error("can't create admin client", zap.Error(err))
		return err
	}

	topicSpecifications := []kafka.TopicSpecification{}
	for _, topicName := range p.topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{Topic: topicName, NumPartitions: 1, ReplicationFactor: 1})
	}

	_, err = adminClient.CreateTopics(context.Background(), topicSpecifications)
	if err != nil {
		p.logger.Error("can't register topics", zap.Error(err))

		return err
	}

	return nil
}
