package subscriber

import (
	"context"
	"gorm.io/gorm"

	"social_todo_list/common"
	"social_todo_list/module/item/storage"
	"social_todo_list/pubsub"

	goservice "github.com/haohmaru3000/go_sdk"
)

type HasItemId interface {
	GetItemId() int
}

//func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext, ctx context.Context) {
//	ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
//	db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
//
//	c, _ := ps.Subscribe(ctx, common.TopicUserLikedItem)
//
//	go func() {
//		defer common.Recovery()
//		for msg := range c {
//			data := msg.Data().(HasItemId)
//
//			if err := storage.NewSQLStore(db).IncreaseLikeCount(ctx, data.GetItemId()); err != nil {
//				log.Println(err)
//			}
//		}
//	}()
//}

func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Increase like count after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			// data := message.Data().(HasItemId)

			data := message.Data().(map[string]interface{})
			itemId := data["item_id"].(float64)

			if err := storage.NewSQLStore(db).IncreaseLikeCount(ctx, int(itemId)); err != nil {
				return err
			}

			_ = message.Ack()
			return nil
		},
	}
}
