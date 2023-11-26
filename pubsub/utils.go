package pubsub

import (
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	common_utils "github.com/dispenal/go-common/utils"
)

func SetRetryOrSetDataToDB(config *common_utils.BaseConfig, msg *pubsub.Message, cb func()) {
	if msg.DeliveryAttempt != nil && *msg.DeliveryAttempt <= 4 {
		time.Sleep(5 * time.Second)
		common_utils.LogInfo(fmt.Sprintf("retry message with messageID: %s, orderingKey: %s", msg.ID, msg.OrderingKey))
		msg.Nack()
	}

	if msg.DeliveryAttempt != nil && *msg.DeliveryAttempt > 4 {
		cb()
		common_utils.LogInfo("acknowledged message")
		msg.Ack()
	}
}

func BuildDescErrorMsg(desc string, err error) string {
	return fmt.Sprintf("%s, Error: %s", desc, err.Error())
}
