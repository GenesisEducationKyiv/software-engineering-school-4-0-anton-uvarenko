package producer

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type emailEventPayload struct {
	Email string `json:"email"`
}

var emailTopicName = "email"

func (p *Producer) ProduceEmailEvent(email string) error {
	payload, err := json.Marshal(emailEventPayload{
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
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
