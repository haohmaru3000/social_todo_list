package cmd

import (
	"gorm.io/gorm"
	"log"

	"to_do_list/common"

	goservice "github.com/haohmaru3000/go_sdk"
	"github.com/haohmaru3000/go_sdk/plugin/storage/sdkgorm"
	"github.com/spf13/cobra"
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
