package kafka

import (
	"context"
	"net"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/topics"
)

func newClient(addr net.Addr) *kafka.Client {
	client := &kafka.Client{
		Addr:    addr,
		Timeout: 5 * time.Second,
	}

	return client
}

func (k *Client) ListTopics() []kafka.Topic {
	client := newClient(kafka.TCP(k.cfg.KafkaBrokers...))
	topics, err := topics.List(context.TODO(), client)
	if err != nil {
		common_utils.LogError("list topic error: " + err.Error())
		return nil
	}
	return topics

}

func (k *Client) CreateTopic(topic string, numPart int) error {
	client := newClient(kafka.TCP(k.cfg.KafkaBrokers...))

	_, err := client.CreateTopics(context.TODO(), &kafka.CreateTopicsRequest{
		Addr: kafka.TCP(k.cfg.KafkaBrokers...),
		Topics: []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     numPart,
				ReplicationFactor: k.cfg.KafkaReplicationFactor,
			},
		},
	})
	if err != nil {
		common_utils.LogError("create topic error: " + err.Error())
		return err
	}

	common_utils.LogInfo("topic created: " + topic)
	return nil
}
