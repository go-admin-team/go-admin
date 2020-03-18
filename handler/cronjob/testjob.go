package cronjob

import (
	"github.com/robfig/cron/v3"
	"log"
)

func TestJob(c *cron.Cron) {
	id, err := c.AddFunc("1 * * * *", func() {

		log.Println("Every hour on the one hour")
	})
	if err != nil {
		log.Println(err)
		log.Println("start error")
	} else {
		log.Println("Start Success; ID: %v", id)
	}
}
