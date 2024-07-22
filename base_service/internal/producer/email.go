package producer

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var emailTopicName = "emails"

type subscribedEventPayload struct {
	Email string `json:"email"`
}

func (p *Producer) ProduceSubscribedEvent(email string) error {
	payload, err := json.Marshal(subscribedEventPayload{
		Email: email,
	})
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
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
	}, nil)
	if err != nil {
		return err
	}

	return nil
}

type unsubscribedEventPayload struct {
	Email string `json:"email"`
}

func (p *Producer) ProduceUnsubscribedEvent(email string) error {
	payload, err := json.Marshal(unsubscribedEventPayload{
		Email: email,
	})
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
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
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
