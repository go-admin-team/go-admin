package config

type Ssl struct {
	KeyStr string
	Pem    string
	Enable bool
	Domain string
}

var SslConfig = new(Ssl)
