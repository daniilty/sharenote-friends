package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

// Consumer - kafka consumer.
type Consumer interface {
	UnmarshalMessage(context.Context, interface{}) (CommitFunc, error)
	Close() error
}

// ConsumerImpl - consumer implementation.
type ConsumerImpl struct {
	reader *kafka.Reader
}

// NewComcumerImpl - ConsumerImpl constructor.
func NewConsumerImpl(topic string, brokers []string, groupID string) *ConsumerImpl {
	return &ConsumerImpl{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
			Dialer: &kafka.Dialer{
				Timeout:   10 * time.Second,
				DualStack: true,
			},
		}),
	}
}

// UnmarshalMessage - fetch message from kafka broker.
func (c *ConsumerImpl) UnmarshalMessage(ctx context.Context, msg interface{}) (CommitFunc, error) {
	kafkaMsg, err := c.reader.FetchMessage(ctx)
	if err != nil {
		return func(ctx2 context.Context) error { return nil }, err
	}

	err = json.Unmarshal(kafkaMsg.Value, &msg)
	if err != nil {
		return func(ctx2 context.Context) error { return nil }, err
	}

	return func(ctx2 context.Context) error {
		return c.reader.CommitMessages(ctx2, kafkaMsg)
	}, nil
}

func (c *ConsumerImpl) Close() error {
	return c.reader.Close()
}
