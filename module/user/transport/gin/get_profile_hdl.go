package ginuser

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"social_todo_list/common"
)

func Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
