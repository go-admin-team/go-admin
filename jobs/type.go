package jobs

type Job interface {
	Run()
	addJob() (int, error)
}

type JobsExec interface {
	Exec()
}

func CallExec(e JobsExec) {
	e.Exec()
}
