package ginitem

import (
	"gorm.io/gorm"
	"net/http"

	"to_do_list/common"
	"to_do_list/module/item/biz"
	"to_do_list/module/item/model"
	"to_do_list/module/item/repository"
	"to_do_list/module/item/storage"
	usrLikeStore "to_do_list/module/userlikeitem/storage"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func ListItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

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

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		store := storage.NewSQLStore(db)
		likeStore := usrLikeStore.NewSQLStore(db)
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		for i := range result {
			result[i].Mask()
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}
