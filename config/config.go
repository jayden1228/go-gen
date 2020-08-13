package config

import (
	"fmt"

	"github.com/timest/env"
)

//EnvConfig EnvConfig
var EnvConfig *config

type config struct {
	ProjectEnv string `env:"PROJECT_ENV" default:"dev"`
	APIVersion string `env:"API_VERSION" default:"Commit ID"`
	Mysql      struct {
		Host    string `default:"127.0.0.1"`
		Port    string `default:"3306"`
		DBName  string `default:"information_schema"`
		User    string `default:"root"`
		Pwd     string `default:"123"`
		Charset string `default:"utf8mb4"`
	}
}

func init() {
	EnvConfig = new(config)
	env.IgnorePrefix()
	err := env.Fill(EnvConfig)
	fmt.Println(EnvConfig)
	if err != nil {
		panic(err)
	}
}
