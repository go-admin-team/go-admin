package cronjob

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func TestJob(c *cron.Cron) {
	id, err := c.AddFunc("1 * * * *", func() {

		log.Println("Every hour on the one hour")
	})
	if err != nil {
		log.Println(err)
		log.Println("start error")
	} else {
		log.Printf("Start Success; ID: %v \r\n", id)
	}
}
