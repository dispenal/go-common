package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	common_utils "github.com/dispenal/go-common/utils"
)

type GooglePubSub interface {
	CreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error)
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	Topic(id string) *pubsub.Topic
	Subscription(id string) *pubsub.Subscription
	Close() error
}

type PubSubClientImpl struct {
	config *common_utils.BaseConfig
	pubSub GooglePubSub
}

func NewGooglePubSub(config *common_utils.BaseConfig) (c *pubsub.Client, err error) {
	ctx := context.Background()
	projectId := config.GcpProjectId

	return pubsub.NewClient(ctx, projectId)
}

func NewPubSubClient(config *common_utils.BaseConfig, pubSub GooglePubSub) PubSubClient {
	return &PubSubClientImpl{config: config, pubSub: pubSub}
}

func (p *PubSubClientImpl) CreateTopicIfNotExists(ctx context.Context, topicName string) (*pubsub.Topic, error) {
	tpc := p.pubSub.Topic(topicName)

	ok, err := tpc.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return tpc, nil
	}

	tpc, err = p.pubSub.CreateTopic(ctx, topicName)
	tpc.EnableMessageOrdering = true

	return tpc, err
}

func (p *PubSubClientImpl) CreateSubscriptionIfNotExists(ctx context.Context, id string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	sub := p.pubSub.Subscription(id)

	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if ok {
		return sub, nil
	}

	return p.pubSub.CreateSubscription(ctx, id, pubsub.SubscriptionConfig{
		Topic:                 topic,
		EnableMessageOrdering: true,
		AckDeadline:           30 * time.Second,
		DeadLetterPolicy: &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     p.config.PubsubDlq,
			MaxDeliveryAttempts: 5,
		},
	})
}

func (p *PubSubClientImpl) PublishTopics(ctx context.Context, topics []*pubsub.Topic, data any, orderingKey string) error {
	var results []*pubsub.PublishResult

	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		topic.EnableMessageOrdering = true
		res := topic.Publish(ctx, &pubsub.Message{
			Data:        message,
			OrderingKey: orderingKey,
		})
		results = append(results, res)
	}

	for _, result := range results {
		id, err := result.Get(ctx)
		if err != nil {
			return err
		}
		common_utils.LogInfo(fmt.Sprintf("publish message with ID: %s", id))
	}

	return nil
}

func (p *PubSubClientImpl) PullMessages(ctx context.Context, id string, topic *pubsub.Topic, callback func(ctx context.Context, msg *pubsub.Message)) error {
	defer p.pubSub.Close()

	sub, err := p.CreateSubscriptionIfNotExists(ctx, id, topic)
	if err != nil {
		return err
	}

	// sub.ReceiveSettings.Synchronous = true
	// sub.ReceiveSettings.MaxOutstandingMessages = 30
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		common_utils.LogInfo(fmt.Sprintf("received message with messageID: %s, topicID: %s, ordering key: %s \n", msg.ID, topic.ID(), msg.OrderingKey))

		callback(ctx, msg)
	})
}

func (p *PubSubClientImpl) Close() error {
	return p.pubSub.Close()
}

func (p *PubSubClientImpl) CheckTopicAndPublish(
	ctx context.Context,
	topicsName []string,
	orderingKey string,
	data any,
) {
	if len(topicsName) == 0 {
		return
	}
	topics := make([]*pubsub.Topic, len(topicsName))
	for i := range topicsName {
		topic, err := p.CreateTopicIfNotExists(ctx, topicsName[i])
		if err != nil {
			common_utils.LogError(fmt.Sprintf("error when creating topic: [%s] \n", topicsName[i]))
			return
		}
		topics[i] = topic
	}
	err := p.PublishTopics(ctx, topics, data, orderingKey)
	common_utils.LogIfError(err)
}
