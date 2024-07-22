package consumer

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer     *kafka.Consumer
	topics       []string
	rateHandler  handler
	emailHandler handler
}

type handler interface {
	Handle(msg []byte) error
}

func NewConsumer(
	rateHandler handler,
	emailHandler handler,
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
		topics:       []string{"emails", "rate"},
		rateHandler:  rateHandler,
		emailHandler: emailHandler,
	}
}

func (c Consumer) InitializeTopics() {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9094",
	})
	if err != nil {
		panic(err)
	}

	topicSpecifications := []kafka.TopicSpecification{}
	for _, topicName := range c.topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{Topic: topicName, NumPartitions: 1, ReplicationFactor: 1})
	}
	_, err = adminClient.CreateTopics(context.Background(), topicSpecifications)
	if err != nil {
		panic(err)
	}
}

func (c Consumer) StartPolling() {
	fmt.Println("start consuming messages")
	c.consumer.SubscribeTopics(c.topics, nil)

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			fmt.Printf("can't read message: %v", err)
			continue
		}

		fmt.Printf("consumed  message: %v", string(msg.Value))
		chosenHandler := c.chooseHandler(msg)
		go func(handler handler, msg *kafka.Message) {
			fmt.Println("handling msg")

			err := handler.Handle(msg.Value)
			if err != nil {
				return
			}

			c.consumer.CommitMessage(msg)
		}(chosenHandler, msg)

	}
}

func (c Consumer) chooseHandler(msg *kafka.Message) handler {
	switch *msg.TopicPartition.Topic {
	case "emails":
		return c.emailHandler

	case "rate":
		return c.rateHandler
	}

	return nil
}
