package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const configFile = "config.yaml"

// Config Configuration struct
type Config struct {
	Application struct {
		HttP struct {
			Port string `yaml:"port"`
		} `yaml:"http"`
	} `yaml:"application"`
}

func GetConfig() (Config, error) {
	var c = Config{}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
