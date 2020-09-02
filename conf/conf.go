package conf

import "github.com/jasontconnell/conf"

type Config struct {
	ConnectionString string `json:"connectionString"`
}

func LoadConfig(filename string) Config {
	cfg := Config{}
	conf.LoadConfig(filename, &cfg)
	return cfg
}
