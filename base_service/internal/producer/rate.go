package producer

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type rateEventPayload struct {
	Rate      float32   `json:"rate"`
	UpdatedAt time.Time `json:"updated_at"`
}

var rateEventName = "rate"

func (p *Producer) ProduceRateEvent(rate float32) error {
	payload, err := json.Marshal(rateEventPayload{
		Rate:      rate,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &rateEventName,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, nil)
	if err != nil {
		return err
	}

	return nil
}
