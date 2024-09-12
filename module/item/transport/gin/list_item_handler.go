package ginitem

import (
	"net/http"
	"to_do_list/common"
	"to_do_list/module/item/biz"
	"to_do_list/module/item/model"
	"to_do_list/module/item/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		queryString.Paging.Process()

		store := storage.NewSQLStore(db)
		business := biz.NewListItemBiz(store)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
