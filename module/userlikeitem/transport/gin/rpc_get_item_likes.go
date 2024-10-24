package ginuserlikeitem

import (
	"gorm.io/gorm"
	"net/http"

	"social_todo_list/common"
	"social_todo_list/module/userlikeitem/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func GetItemLikes(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		type RequestData struct {
			Ids []int `json:"ids"`
		}

		var data RequestData

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)

		mapRs, err := store.GetItemLikes(c.Request.Context(), data.Ids)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(mapRs))
	}
}
