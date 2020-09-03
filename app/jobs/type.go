package jobs

type Job interface {
	Run()
	addJob() (int, error)
}

type JobsExec interface {
	Exec(arg interface{}) error
}

func CallExec(e JobsExec, arg interface{}) error {
	return e.Exec(arg)
}
