package ginuser

import (
	"gorm.io/gorm"
	"net/http"

	"to_do_list/common"
	"to_do_list/module/user/biz"
	"to_do_list/module/user/model"
	"to_do_list/module/user/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func Register(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		var data model.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()
		biz := biz.NewRegisterBusiness(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}
