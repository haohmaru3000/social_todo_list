package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"to_do_list/common"
	"to_do_list/module/user/biz"
	"to_do_list/module/user/model"
	"to_do_list/module/user/storage"
)

func Register(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
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
