package configuration

import (
	"github.com/KaiserWerk/Maestro/internal/global"
	"io/ioutil"
	"os"

	"github.com/KaiserWerk/Maestro/internal/assets"
	"github.com/KaiserWerk/Maestro/internal/entity"

	"gopkg.in/yaml.v2"
)

var (
	configFile = "app.yaml"
	appConfig *entity.AppConfig
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

	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("could not read configuration file: " + err.Error())
	}

	var conf entity.AppConfig
	if err := yaml.Unmarshal(content, &conf); err != nil {
		panic("could not unmarshal configuration: " + err.Error())
	}

	if created {
		if e := os.Getenv(global.EnvBindAddress); e != "" {
			conf.App.BindAddress = e
		}
		if e := os.Getenv(global.EnvAuthToken); e != "" {
			conf.App.AuthToken = e
		}
		if e := os.Getenv(global.EnvCertFile); e != "" {
			conf.App.CertificateFile = e
		}
		if e := os.Getenv(global.EnvKeyFile); e != "" {
			conf.App.KeyFile = e
		}
		if e := os.Getenv(global.EnvDatabaseDriver); e != "" {
			conf.Database.Driver = e
		}
		if e := os.Getenv(global.EnvDatabaseDSN); e != "" {
			conf.Database.DSN = e
		}
	}

	appConfig = &conf

	return appConfig, created, nil
}

func GetConfiguration() *entity.AppConfig {
	return appConfig
}
