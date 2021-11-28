package configuration

import (
	"io/ioutil"
	"os"

	"github.com/KaiserWerk/Maestro/internal/assets"
	"github.com/KaiserWerk/Maestro/internal/entity"

	"gopkg.in/yaml.v2"
)

var (
	configFile = "app.yaml"
)

func SetFile(file string) {
	configFile = file
}

func Setup() (*entity.AppConfig, bool, error) {
	var created bool
	if _, err := os.Stat(configFile); os.IsNotExist(err)  {
		content, err := assets.ReadConfigurationFile("app.dist.yaml")
		if err != nil {
			return nil, created, err
		}

		err = ioutil.WriteFile(configFile, content, 0744)
		if err != nil {
			return nil, created, err
		}
		created = true
	}

	conf, err := GetConfiguration()
	return conf, created, err
}

func GetConfiguration() (*entity.AppConfig, error) {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var conf entity.AppConfig
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
