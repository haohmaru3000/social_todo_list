package subscriber

import (
	"context"
	"gorm.io/gorm"

	"social_todo_list/common"
	"social_todo_list/module/item/storage"
	"social_todo_list/pubsub"

	goservice "github.com/haohmaru3000/go_sdk"
)

func DecreaseLikeCountAfterUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Decrease like count after user unlikes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			// data := message.Data().(HasItemId)
			data := message.Data().(map[string]interface{})
			itemId := data["item_id"].(float64)

			if err := storage.NewSQLStore(db).DecreaseLikeCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = message.Ack()
			return nil
		},
	}
}
