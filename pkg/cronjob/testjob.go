package cronjob

import (
	"github.com/robfig/cron/v3"
	"log"
	tools2 "go-admin/tools"
)

func TestJob(c *cron.Cron) {
	id, err := c.AddFunc("1 * * * *", func() {

		tools2.Logger.Println("Every hour on the one hour")
	})
	if err != nil {
		tools2.Logger.Println(err)
		log.Println("start error")
	} else {
		log.Printf("Start Success; ID: %v \r\n", id)
	}
}
