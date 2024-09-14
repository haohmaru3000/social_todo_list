package cmd

import (
	"log"

	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/spf13/cobra"
	"gorm.io/gorm"

	"to_do_list/common"
	"to_do_list/plugin/sdkgorm"
)

var cronDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run demo cron job",
	Run: func(cmd *cobra.Command, args []string) {

		service := goservice.New(
			goservice.WithName("social-todo-list"),
			goservice.WithVersion("1.0.0"),
			goservice.WithInitRunnable(sdkgorm.NewGormDB("main.mysql", common.PluginDBMain)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		db := service.MustGet(common.PluginDBMain).(*gorm.DB)

		log.Println("I am demo cron with DB connection:", db)

	},
}
