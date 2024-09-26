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
	"github.com/haohmaru3000/go_sdk/plugin/tokenprovider"
)

func Login(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		tokenProvider := serviceCtx.MustGet(common.PluginJWT).(tokenprovider.Provider)

		var loginUserData model.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := storage.NewSQLStore(db)
		md5 := common.NewMd5Hash()

		business := biz.NewLoginBusiness(store, tokenProvider, md5, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
