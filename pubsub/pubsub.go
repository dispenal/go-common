package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type PubSubClient interface {
	CreateTopicIfNotExists(ctx context.Context, topicName string) (*pubsub.Topic, error)
	CreateSubscriptionIfNotExists(ctx context.Context, id string, topic *pubsub.Topic) (*pubsub.Subscription, error)
	PublishTopics(ctx context.Context, topics []*pubsub.Topic, data any, orderingKey string) error
	PullMessages(ctx context.Context, id string, topic *pubsub.Topic, callback func(ctx context.Context, msg *pubsub.Message)) error
	Close() error
	CheckTopicAndPublish(ctx context.Context, topicsName []string, orderingKey string, data any)
}
