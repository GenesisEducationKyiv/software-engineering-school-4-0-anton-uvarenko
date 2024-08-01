package consumer

import (
	"context"
	"fmt"

	"github.com/VictoriaMetrics/metrics"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

var consumedEventsTotal = metrics.NewCounter("consumed_events_total")

type Consumer struct {
	consumer     *kafka.Consumer
	topics       []string
	emailHandler handler
	logger       *zap.Logger
}

type handler interface {
	Handle(msg *kafka.Message) error
}

func NewConsumer(
	emailHandler handler,
	logger *zap.Logger,
) *Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
		"group.id":          "emails",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	return &Consumer{
		consumer:     consumer,
		topics:       []string{"emails"},
		emailHandler: emailHandler,
		logger:       logger.With(zap.String("service", "Consumer")),
	}
}

func (c Consumer) InitializeTopics() {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		c.logger.Error("can't initialize topics", zap.Error(err))
		panic(err)
	}

	topicSpecifications := []kafka.TopicSpecification{}
	for _, topicName := range c.topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{Topic: topicName, NumPartitions: 1, ReplicationFactor: 1})
	}
	_, err = adminClient.CreateTopics(context.Background(), topicSpecifications)
	if err != nil {
		c.logger.Error("can't create topics", zap.Error(err))
		panic(err)
	}
}

func (c Consumer) StartPolling() {
	logger := c.logger.With(zap.String("method", "StartPolling"))

	fmt.Println("start consuming messages")
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		logger.Error("can't subscribe to topics", zap.Error(err))
		return
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logger.Error("can't read message", zap.Error(err))
			continue
		}

		consumedEventsTotal.Inc()

		logger.Info("consumed  message", zap.Any("message", msg))
		chosenHandler := c.chooseHandler(msg)
		go func(handler handler, msg *kafka.Message, logger *zap.Logger) {
			err := handler.Handle(msg)
			if err != nil {
				logger.Error("can't handler message", zap.Error(err))
				return
			}

			_, err = c.consumer.CommitMessage(msg)
			if err != nil {
				logger.Error("can't commit message", zap.Error(err))
			}
		}(chosenHandler, msg, logger)

	}
}

func (c Consumer) chooseHandler(msg *kafka.Message) handler {
	if *msg.TopicPartition.Topic == "emails" {
		return c.emailHandler
	}

	return nil
}
