package event

import (
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type KafkaAdapter interface {
	SendEvent(event *Event) error
	SendEventWithPayload(eventType EventType, payload any) error
	Close() error
}

type kafkaAdapter struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaAdapter creates a new KafkaAdapter.
func NewKafkaAdapter(brokers []string, topic string) (KafkaAdapter, error) {
	log.Infof("creating kafka adapter with brokers:%v and topic %v", brokers, topic)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %v", err)
	}

	return &kafkaAdapter{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *kafkaAdapter) SendEvent(event *Event) error {
	message := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.StringEncoder(event.String()),
	}

	_, _, err := kp.producer.SendMessage(message)
	if err != nil {
		log.Errorf("event:%v sending error:%v", event, err)
		return err
	}

	return nil
}

func (kp *kafkaAdapter) SendEventWithPayload(eventType EventType, payload any) error {
	event, err := NewEvent(eventType, payload)
	if err != nil {
		return err
	}
	err = kp.SendEvent(event)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the KafkaAdapter.
func (kp *kafkaAdapter) Close() error {
	return kp.producer.Close()
}
