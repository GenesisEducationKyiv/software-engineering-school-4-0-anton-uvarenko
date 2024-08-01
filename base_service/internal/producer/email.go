package producer

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.uber.org/zap"
)

var emailTopicName = "emails"

type subscribedEventPayload struct {
	Email string `json:"email"`
}

func (p *Producer) ProduceSubscribedEvent(email string) error {
	logger := p.logger.With(zap.String("method", "ProduceSubscribedEvent"))

	payload, err := json.Marshal(subscribedEventPayload{
		Email: email,
	})
	if err != nil {
		logger.Error("can't unmarshal payload", zap.Error(err))
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &emailTopicName,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
		Headers: []kafka.Header{
			{
				Key:   "origin",
				Value: []byte("user_subscribed"),
			},
		},
	})
	if err != nil {
		logger.Error("can't produce event", zap.Error(err))
		return err
	}

	return nil
}

type unsubscribedEventPayload struct {
	Email string `json:"email"`
}

func (p *Producer) ProduceUnsubscribedEvent(email string) error {
	logger := p.logger.With(zap.String("method", "ProduceUnsubscribedEvent"))

	payload, err := json.Marshal(unsubscribedEventPayload{
		Email: email,
	})
	if err != nil {
		logger.Error("can't unmarshal payload", zap.Error(err))
		return err
	}

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &emailTopicName,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
		Headers: []kafka.Header{
			{
				Key:   "origin",
				Value: []byte("user_unsubscribed"),
			},
		},
	})
	if err != nil {
		logger.Error("can't produce event", zap.Error(err))
		return err
	}

	return nil
}
