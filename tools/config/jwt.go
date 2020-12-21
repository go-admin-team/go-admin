package config

type Jwt struct {
	Secret  string
	Timeout int64
}

var JwtConfig = new(Jwt)
