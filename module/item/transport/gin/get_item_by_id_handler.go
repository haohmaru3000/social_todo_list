package ginitem

import (
	"gorm.io/gorm"
	"net/http"

	"to_do_list/common"
	"to_do_list/module/item/biz"
	"to_do_list/module/item/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func GetItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		business := biz.NewGetItemBiz(store)

		data, err := business.GetItemById(c.Request.Context(), int(uid.GetLocalID()))

		if err != nil {
			panic(err)
		}

		data.Mask()

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}
