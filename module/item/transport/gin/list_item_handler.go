package ginitem

import (
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"net/http"

	"social_todo_list/common"
	"social_todo_list/demogrpc/demo"
	"social_todo_list/module/item/biz"
	"social_todo_list/module/item/model"
	"social_todo_list/module/item/repository"
	"social_todo_list/module/item/storage"
	"social_todo_list/module/item/storage/rpc"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
)

func ListItem(serviceCtx goservice.ServiceContext, client demo.ItemLikeServiceClient) func(*gin.Context) {
	return func(c *gin.Context) {
		propagator := otel.GetTextMapPropagator()
		extractCtx := propagator.Extract(c, propagation.HeaderCarrier(c.Request.Header))

		_, span := otel.Tracer("Item.Transport").Start(extractCtx, "item.transport.list")
		defer span.End()

		traceParent := c.Request.Header.Get("traceparent")
		fmt.Println("traceParent: " + traceParent)

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		//apiItemCaller := serviceCtx.MustGet(common.PluginItemAPI).(interface {
		//	GetServiceURL() string
		//})

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
		//likeStore := restapi.New(apiItemCaller.GetServiceURL(), serviceCtx.Logger("restapi.itemlikes"))

		likeStore := rpc.NewClient(client)
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
