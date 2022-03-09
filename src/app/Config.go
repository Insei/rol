package app

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"runtime"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		SqlConnectionString string `yaml:"sqlConnectionString"`
		DbName              string `yaml:"dbname"`
		Username            string `yaml:"user"`
		Password            string `yaml:"password"`
	} `yaml:"database"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}

func readConfig() *Config {
	cfg := &Config{}
	_, filename, _, _ := runtime.Caller(1)
	f, err := os.Open(path.Join(path.Dir(filename), "../appConfig.yml"))
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
	return cfg
}

func GetConfig() *Config {
	cfg := readConfig()
	return cfg
}
