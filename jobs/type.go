package jobs

type Job interface {
	Run()
	addJob() (int, error)
}
