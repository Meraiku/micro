package notification

import (
	"context"
	"log"

	"github.com/meraiku/micro/notification/pkg/consumer"
)

var (
	brokers = []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"}
)

type Service struct {
	recieve chan string
}

func New(ctx context.Context) (*Service, error) {

	topic := "user"

	notifChan, err := consumer.NewGroup(ctx, brokers, "notification", topic)
	if err != nil {
		return nil, err
	}

	return &Service{
		recieve: notifChan,
	}, nil
}

func (s *Service) Read() {
	for {
		select {
		case msg := <-s.recieve:
			log.Printf("got message in notification service: %v", msg)
		}
	}
}
