package cmd

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"

	"to_do_list/common"
	"to_do_list/memcache"
	"to_do_list/middleware"
	ginitem "to_do_list/module/item/transport/gin"
	"to_do_list/module/upload"
	userstorage "to_do_list/module/user/storage"
	ginuser "to_do_list/module/user/transport/gin"
	ginuserlikeitem "to_do_list/module/userlikeitem/transport/gin"
	"to_do_list/plugin/simple"
	"to_do_list/pubsub"
	"to_do_list/subscriber"

	"github.com/gin-gonic/gin"
	goservice "github.com/haohmaru3000/go_sdk"
	"github.com/haohmaru3000/go_sdk/plugin/jaeger"
	"github.com/haohmaru3000/go_sdk/plugin/rpccaller"
	"github.com/haohmaru3000/go_sdk/plugin/storage/sdkgorm"
	"github.com/haohmaru3000/go_sdk/plugin/tokenprovider/jwt"
	"github.com/spf13/cobra"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginItemAPI)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
		goservice.WithInitRunnable(jaeger.NewJaeger(common.PluginTracingService)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			// Example for Simple Plugin
			type CanGetValue interface {
				GetValue() string
			}
			log.Println(service.MustGet("simple").(CanGetValue).GetValue())
			/////////

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			authCache := memcache.NewUserCaching(memcache.NewCaching(), authStore)

			middlewareAuth := middleware.RequiredAuth(authCache, service)

			v1 := engine.Group("/v1")
			{
				v1.PUT("/upload", upload.Upload(service))

				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(service))
					items.GET("", ginitem.ListItem(service))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))

					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
				}
			}

			rpc := v1.Group("rpc")
			{
				rpc.POST("/get_item_likes", ginuserlikeitem.GetItemLikes(service))
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		/***** Start TRACING *****/
		type TracingService interface {
			Run() error
		}

		je := service.MustGet("jaeger").(TracingService)
		if err := je.Run(); err != nil {
			log.Fatalln(err)
		}
		/***** End TRACING *****/

		_ = subscriber.NewEngine(service).Start()

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronDemoCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
