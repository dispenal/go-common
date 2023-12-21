package kafka

import (
	"context"
	"testing"

	common_utils "github.com/dispenal/go-common/utils"
)

func TestKafkaConenection(t *testing.T) {
	config, err := common_utils.LoadBaseConfig("../", "test")
	if err != nil {
		t.Error(err)
	}
	kafkaClient := NewKafkaClient(config)
	err = kafkaClient.NewPublisher()
	if err != nil {
		t.Error(err)
	}

	kafkaClient.NewConsumer()

	err = kafkaClient.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestKafkaPublisher(t *testing.T) {
	config, err := common_utils.LoadBaseConfig("../", "test")
	if err != nil {
		t.Error(err)
	}
	kafkaClient := NewKafkaClient(config)
	err = kafkaClient.NewPublisher()
	if err != nil {
		t.Error(err)
	}

	err = kafkaClient.CreateTopic("tester", -1)
	if err != nil {
		t.Error(err)
	}

	event := NewEvent("test", []byte("test"))

	err = kafkaClient.Publish(context.Background(), "tester", *event)
	if err != nil {
		t.Error(err)
	}

	err = kafkaClient.Close()
	if err != nil {
		t.Error(err)
	}
}
