package consumer

import (
	"context"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func NewGroup(
	ctx context.Context,
	brokers []string,
	groupID, topic string,
) (chan *ConsumerMessage, error) {
	config := sarama.NewConfig()

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 100 * time.Millisecond

	config.Consumer.MaxWaitTime = 500 * time.Millisecond

	config.Consumer.Return.Errors = true

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	notifChan, err := subscribeGroup(ctx, topic, consumer)
	if err != nil {
		return nil, err
	}

	return notifChan, nil
}

func subscribeGroup(
	ctx context.Context,
	topic string,
	consumer sarama.ConsumerGroup,
) (chan *ConsumerMessage, error) {

	handler, notifChan := NewConsumer()

	go func() {
		defer consumer.Close()

		for {
			err := consumer.Consume(ctx, []string{topic}, handler)
			if err != nil {
				log.Printf("Error from consumer: %v", err)
				log.Println("Closing consumer...")
				return
			}
		}
	}()

	return notifChan, nil
}

type Consumer struct {
	notif chan *ConsumerMessage
}

func NewConsumer() (*Consumer, chan *ConsumerMessage) {
	notifChan := make(chan *ConsumerMessage)
	return &Consumer{
		notif: notifChan,
	}, notifChan
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer Group are been rebalanced")
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer Group will be rebalanced soon!")
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {

		c.notif <- msg

		session.MarkMessage(msg, "")
	}

	return nil
}
