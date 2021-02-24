package config

import "gorm.io/plugin/dbresolver"

type Database struct {
	Driver          string
	Source          string
	ConnMaxIdleTime int
	ConnMaxLifetime int
	MaxIdleConns    int
	MaxOpenConns    int
	Registers       []DBResolverConfig
}

type DBResolverConfig struct {
	Sources  []string
	Replicas []string
	Policy   string
	Tables   []string
}

var (
	DatabaseConfig  = new(Database)
	DatabasesConfig = make(map[string]*Database)
	Policies        = map[string]dbresolver.Policy{
		"random": dbresolver.RandomPolicy{},
	}
)
