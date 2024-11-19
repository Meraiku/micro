package producer

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer AsyncProducer
	topic    string
	msgs     chan *Message
}

func New(brokers []string, topic string) (*Producer, error) {
	producer, err := newAsync(brokers)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
		msgs:     make(chan *Message, 1),
	}, nil
}

func (p *Producer) Run() {
	go func() {
		for range p.producer.Successes() {
			log.Println("Successfully produced message")
		}
	}()

	go func() {
		for err := range p.producer.Errors() {
			log.Printf("Failed to produce message: %v\n", err)
		}
	}()

	go func() {
		for msg := range p.msgs {
			p.producer.Input() <- msg
		}
	}()
}

func (p *Producer) Close() {
	if err := p.producer.Close(); err != nil {
		log.Println(err)
	}
}

func (p *Producer) Send(key string, value string) {
	message := []byte(value)

	msg := PrepareMessage(p.topic, key, message)

	p.msgs <- msg
}

func newAsync(brokers []string) (AsyncProducer, error) {

	config := sarama.NewConfig()

	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.MaxOpenRequests = 1

	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func PrepareMessage(topic, key string, message []byte) *Message {
	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(message),
		Partition: -1,
	}

	return msg
}
