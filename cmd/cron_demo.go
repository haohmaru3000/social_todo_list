package cmd

import (
	"gorm.io/gorm"
	"log"

	"to_do_list/common"

	goservice "github.com/haohmaru3000/go_sdk"
	"github.com/haohmaru3000/go_sdk/plugin/storage/sdkgorm"
	"github.com/spf13/cobra"
)

// Ví dụ cho cronjob update lại số like count trên toàn bộ table todo_items

// UPDATE todo_items ti INNER JOIN (
// SELECT item_id, COUNT(item_id) as `count` FROM `user_like_items`
// GROUP BY item_id
// ) c ON c.item_id = ti.id SET ti.liked_count = c.count

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
