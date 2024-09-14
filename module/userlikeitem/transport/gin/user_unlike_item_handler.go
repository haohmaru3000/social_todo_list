package ginuserlikeitem

import (
	"net/http"

	"to_do_list/common"
	itemStorage "to_do_list/module/item/storage"
	"to_do_list/module/userlikeitem/biz"
	"to_do_list/module/userlikeitem/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
	"gorm.io/gorm"
)

func UnlikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		itemStore := itemStorage.NewSQLStore(db)
		business := biz.NewUserUnlikeItemBiz(store, itemStore)

		if err := business.UnlikeItem(c.Request.Context(), requester.GetUserId(), int(id.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
