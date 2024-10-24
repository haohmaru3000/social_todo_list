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

			data := message.Data().(HasItemId)

			return storage.NewSQLStore(db).DecreaseLikeCount(ctx, data.GetItemId())
		},
	}
}
