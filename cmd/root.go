package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"

	"social_todo_list/common"
	"social_todo_list/demogrpc/demo"
	"social_todo_list/memcache"
	"social_todo_list/middleware"
	ginitem "social_todo_list/module/item/transport/gin"
	"social_todo_list/module/upload"
	userstorage "social_todo_list/module/user/storage"
	ginuser "social_todo_list/module/user/transport/gin"
	"social_todo_list/module/userlikeitem/storage"
	ginuserlikeitem "social_todo_list/module/userlikeitem/transport/gin"
	"social_todo_list/module/userlikeitem/transport/rpc"
	"social_todo_list/plugin/simple"
	// "social_todo_list/pubsub"
	"social_todo_list/plugin/nats"
	"social_todo_list/subscriber"

	goservice "github.com/haohmaru3000/go_sdk"
	"github.com/haohmaru3000/go_sdk/plugin/jaeger"
	"github.com/haohmaru3000/go_sdk/plugin/rpccaller"
	"github.com/haohmaru3000/go_sdk/plugin/storage/sdkgorm"
	"github.com/haohmaru3000/go_sdk/plugin/storage/sdkredis"
	"github.com/haohmaru3000/go_sdk/plugin/tokenprovider/jwt"
	"github.com/spf13/cobra"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		// goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(nats.NewNATS(common.PluginPubSub)),
		goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginItemAPI)),
		goservice.WithInitRunnable(jaeger.NewJaeger(common.PluginTracingService)),
		goservice.WithInitRunnable(sdkredis.NewRedisDB("redis", common.PluginRedis)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
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

		/***** Setup gRPC *****/
		// gRPC Server
		address := "0.0.0.0:50051"
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		fmt.Printf("Server is listening on %v ...\n", address)
		s := grpc.NewServer()
		db := service.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewSQLStore(db)
		demo.RegisterItemLikeServiceServer(s, rpc.NewRPCService(store))
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalln(err)
			}
		}()

		// gRPC Client
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		cc, err := grpc.NewClient("localhost:50051", opts)
		if err != nil {
			log.Fatal(err)
		}
		client := demo.NewItemLikeServiceClient(cc)
		////////

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
			authCache := memcache.NewUserCaching(memcache.NewRedisCache(service), authStore)

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
					items.GET("", ginitem.ListItem(service, client))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))
					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
				}

				rpc := v1.Group("rpc")
				{
					rpc.POST("/get_item_likes", ginuserlikeitem.GetItemLikes(service))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		/***** Setup TRACING *****/
		type TracingService interface {
			Run() error
		}
		je := service.MustGet("jaeger").(TracingService)
		if err := je.Run(); err != nil {
			log.Fatalln(err)
		}

		/***** Setup PUBSUB *****/
		_ = subscriber.NewEngine(service).Start()

		/***** Start Service *****/
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
