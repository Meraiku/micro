package notification

import (
	"context"

	"github.com/meraiku/micro/notification/pkg/consumer"
	"github.com/meraiku/micro/notification/pkg/metrics"
	"github.com/meraiku/micro/pkg/logging"
)

var (
	brokers = []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"}
)

type Service struct {
	recieve chan *consumer.ConsumerMessage
	notif   *metrics.NotifMetric
}

func New(ctx context.Context, notif *metrics.NotifMetric) (*Service, error) {

	topic := "user"

	notifChan, err := consumer.NewGroup(ctx, brokers, "notification", topic)
	if err != nil {
		return nil, err
	}

	return &Service{
		recieve: notifChan,
		notif:   notif,
	}, nil
}

func (s *Service) Read(ctx context.Context) {
	log := logging.L(ctx)
	for {
		select {
		case msg := <-s.recieve:
			log.Debug(
				"incrementing notification counter",
			)

			go s.notif.IncrNotifications()

			log.Debug(
				"message received",
				logging.String("key", string(msg.Key)),
				logging.String("message", string(msg.Value)),
			)

		// Process msg

		// Send somewhere

		case <-ctx.Done():
			log.Info("Stopping reading notifications...")
			return
		}
	}
}
