package ginitem

import (
	"gorm.io/gorm"
	"net/http"

	"to_do_list/common"
	"to_do_list/module/item/biz"
	"to_do_list/module/item/model"
	"to_do_list/module/item/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func CreateItem(serviceCtx goservice.ServiceContext) func(ctx *gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		var itemData model.TodoItemCreation

		if err := c.ShouldBind(&itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		itemData.UserId = requester.GetUserId()

		store := storage.NewSQLStore(db)
		business := biz.NewCreateItemBiz(store)

		if err := business.CreateNewItem(c.Request.Context(), &itemData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(itemData.Id))
	}
}
