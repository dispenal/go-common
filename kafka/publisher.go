package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dispenal/go-common/tracer"
	common_utils "github.com/dispenal/go-common/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

func (k *Client) IsReaderConnected() bool {
	return len(k.readers) != 0
}
func (k *Client) NewPublisher() error {
	if len(k.cfg.KafkaBrokers) == 0 {
		return errors.New("not found broker")
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(k.cfg.KafkaBrokers...),
		Balancer:               &kafka.RoundRobin{},
		BatchTimeout:           15 * time.Millisecond,
		AllowAutoTopicCreation: true,
	}

	common_utils.LogInfo("writer created")
	k.writer = w
	return nil
}

func (k *Client) Publish(ctx context.Context, topic string, msg any) error {
	if !k.IsWriters() {
		return errors.New("writers not created")
	}
	if topic == "" {
		return errors.New("topic not empty")
	}

	dataSender, err := json.Marshal(msg)
	if err != nil {
		return errors.New("message of data sender can not marshal")
	}
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = k.writer.WriteMessages(ctx, kafka.Message{
			Topic: topic,
			Key:   []byte(hashMessage(dataSender)),
			Value: dataSender,
			Headers: []kafka.Header{
				protocol.Header{
					Key:   "origin",
					Value: []byte(k.cfg.ServiceName),
				},
			},
		})

		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break

	}
	return nil
}

func (k *Client) PublishWithTracer(ctx context.Context, topic string, msg any) error {
	if !k.IsWriters() {
		return errors.New("writers not created")
	}
	if topic == "" {
		return errors.New("topic not empty")
	}

	dataSender, err := json.Marshal(msg)
	if err != nil {
		return errors.New("message of data sender can not marshal")
	}
	headers := tracer.GetKafkaTracingHeadersFromSpanCtx(ctx)

	headers = append(headers, kafka.Header{
		Key:   "origin",
		Value: []byte(k.cfg.ServiceName),
	})

	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = k.writer.WriteMessages(ctx, kafka.Message{
			Topic:   topic,
			Key:     []byte(hashMessage(dataSender)),
			Value:   dataSender,
			Headers: headers,
		})

		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			return err
		}
		break

	}
	return nil
}

func (k *Client) publishToDLQ(ctx context.Context, m kafka.Message) error {
	if !k.IsWriters() {
		return errors.New("writers not created")
	}
	if m.Topic == "" {
		return errors.New("topic not empty")
	}

	m.Topic = k.cfg.KafkaDlqTopic

	m.Headers = append(m.Headers, kafka.Header{
		Key:   "origin",
		Value: []byte(k.cfg.ServiceName),
	})

	err := k.writer.WriteMessages(ctx, m)
	return err
}