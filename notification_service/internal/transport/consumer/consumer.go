package consumer

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer     *kafka.Consumer
	topics       []string
	emailHandler handler
}

type handler interface {
	Handle(msg *kafka.Message) error
}

func NewConsumer(
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
		topics:       []string{"emails"},
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
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		fmt.Printf("can't SubscribeTopics: %v", err)
		return
	}

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

			err := handler.Handle(msg)
			if err != nil {
				return
			}

			_, err = c.consumer.CommitMessage(msg)
			if err != nil {
				fmt.Printf("can't commit message: %v", err)
			}
		}(chosenHandler, msg)

	}
}

func (c Consumer) chooseHandler(msg *kafka.Message) handler {
	if *msg.TopicPartition.Topic == "emails" {
		return c.emailHandler
	}

	return nil
}
