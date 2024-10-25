package subscriber

import (
	"context"
	"log"

	"social_todo_list/pubsub"

	goservice "github.com/haohmaru3000/go_sdk"
)

type HasUserId interface {
	GetUserId() int
}

func PushNotificationAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Push notification after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			//db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			// data := message.Data().(HasUserId)

			data := message.Data().(map[string]interface{})
			userId := data["user_id"].(float64)
			log.Println("Push notification to user id:", int(userId))

			return nil
		},
	}
}
