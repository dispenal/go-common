package kafka

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dispenal/go-common/tracer"
	common_utils "github.com/dispenal/go-common/utils"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/codes"
)

func (k *Client) NewConsumer() {
	batchSize := int(10e6) // 10MB
	dialer := &kafka.Dialer{
		Timeout:   3 * time.Second,
		DualStack: true,
		KeepAlive: 5 * time.Second,
		ClientID:  RandStringBytes(5),
	}
	k.readers = make(map[string]*kafka.Reader)
	for _, topic := range k.cfg.KafkaTopics {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  k.cfg.KafkaBrokers,
			GroupID:  k.cfg.KafkaGroupID,
			Topic:    topic,
			Dialer:   dialer,
			MaxBytes: batchSize,
		})
		if r == nil {
			common_utils.LogError("empty reader, please check kafka connection")
		}
		common_utils.LogInfo(fmt.Sprintf("Listen: %s, %d, [%s]", r.Stats().Partition, r.Stats().QueueCapacity, r.Stats().Topic))
		k.readers[topic] = r
	}
}

func (k *Client) IsWriters() bool {
	return k.writer != nil
}
func (k *Client) Close() error {
	for _, r := range k.readers {
		r.Close()
	}
	return nil
}

// Listen manual listen
// need call msg.Commit() when process done
// recommend for this process
func (k *Client) Listen(f HandlerFunc) error {
	for _, r := range k.readers {
		ctx := context.Background()

		go func(r *kafka.Reader) {
			for {
				m, err := r.FetchMessage(ctx) // is not auto commit
				if err != nil && errors.Is(err, io.ErrUnexpectedEOF) {
					break
				}
				if err != nil && errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					common_utils.LogIfError(err)
					continue
				}
				retries := 1
				errorMsg := ""

				headers := tracer.TextMapCarrierFromKafkaMessageHeaders(m.Headers)

				msg := &Message{
					Offset:    m.Offset,
					Partition: m.Partition,
					Topic:     m.Topic,
					Headers:   headers,
					Body:      m.Value,
					Timestamp: m.Time.Unix(),
					Key:       string(m.Key),
					Retry:     retries,
					Commit: func() error {
						if err := r.CommitMessages(ctx, m); err != nil {
							return err
						}
						return nil
					},
					MoveToDLQ: func() error {
						return k.publishToDLQ(ctx, m)
					},
				}

				for {
					if retries > k.cfg.KafkaDlqRetry {
						spanCtx, span := tracer.StartAndTraceKafkaConsumer(ctx, headers, "kafkaConsumer.publishToDLQ")
						span.RecordError(errors.New(errorMsg))
						span.SetStatus(codes.Error, errorMsg)
						span.SetAttributes(tracer.BuildAttribute(msg)...)

						common_utils.LogError(fmt.Sprintf("failed process message: %s, will move to DLQ", string(m.Key)))

						m.Headers = append(m.Headers, kafka.Header{
							Key:   "error",
							Value: []byte(errorMsg),
						})

						if err := k.publishToDLQ(spanCtx, m); err != nil {
							tracer.TraceErr(spanCtx, err)
							common_utils.LogError(fmt.Sprintf("failed move message to DLQ: %s", string(m.Key)))
						}

						if err := r.CommitMessages(spanCtx, m); err != nil {
							tracer.TraceErr(spanCtx, err)
							common_utils.LogError(fmt.Sprintf("failed commit message after publish DLQ: %s", string(m.Key)))
						}
						span.End()

						break
					}

					if err := f(ctx, msg); err != nil {
						common_utils.LogError(fmt.Sprintf("failed process message %s with error %v, will retry %d/%d", string(m.Key), err, retries, k.cfg.KafkaDlqRetry))
						errorMsg = err.Error()
						time.Sleep(k.Backoff.NextBackOff())
						retries++
						msg.Retry = retries
						continue
					}
					break
				}

			}
		}(r)
	}
	return nil
}
