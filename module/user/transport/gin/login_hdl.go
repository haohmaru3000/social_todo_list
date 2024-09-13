package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"to_do_list/common"
	"to_do_list/component/tokenprovider"
	"to_do_list/module/user/biz"
	"to_do_list/module/user/model"
	"to_do_list/module/user/storage"
)

func Login(db *gorm.DB, tokenProvider tokenprovider.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
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
