package cmd

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"to_do_list/common"
	"to_do_list/component/tokenprovider/jwt"
	"to_do_list/middleware"
	ginitem "to_do_list/module/item/transport/gin"
	"to_do_list/module/upload"
	userstorage "to_do_list/module/user/storage"
	ginuser "to_do_list/module/user/transport/gin"
	"to_do_list/plugin/sdkgorm"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("social-todo-list"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		systemSecret := os.Getenv("SECRET")

		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStore := userstorage.NewSQLStore(db)
			tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
			middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

			v1 := engine.Group("/v1")
			{
				v1.PUT("/upload", upload.Upload(db))

				v1.POST("/register", ginuser.Register(db))
				v1.POST("/login", ginuser.Login(db, tokenProvider))
				v1.GET("/profile", middlewareAuth, ginuser.Profile())

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(db))
					items.GET("", ginitem.ListItem(db))
					items.GET("/:id", ginitem.GetItem(db))
					items.PATCH("/:id", ginitem.UpdateItem(db))
					items.DELETE("/:id", ginitem.DeleteItem(db))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	// TransAddPoint outenv as sub command
	rootCmd.AddCommand(outEnvCmd)
	// rootCmd.AddCommand(cronjob)

	// rootCmd.AddCommand(startSubUserLikedRestaurantCmd)
	// rootCmd.AddCommand(startSubUserDislikedRestaurantCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
