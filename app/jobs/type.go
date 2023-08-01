package jobs

import "github.com/robfig/cron/v3"

type Job interface {
	Run()
	addJob(*cron.Cron) (int, error)
}

type JobExec interface {
	Exec(arg interface{}) error
}

func CallExec(e JobExec, arg interface{}) error {
	return e.Exec(arg)
}
