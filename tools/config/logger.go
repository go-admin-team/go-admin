package config

type Logger struct {
	Path       string
	Level      string
	Stdout     bool
	EnabledBUS bool
	EnabledREQ bool
	EnabledDB  bool
	EnabledJOB bool `default:"false"`
}

var LoggerConfig = new(Logger)
