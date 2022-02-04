package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"

	"github.com/KaiserWerk/Maestro/internal/assets"
)

type AppConfig struct {
	BindAddress     string        `yaml:"bind_address" envconfig:"BIND_ADDRESS"`
	AuthToken       string        `yaml:"auth_token" envconfig:"AUTH_TOKEN"`
	DieAfter        time.Duration `yaml:"die_after" envconfig:"DIE_AFTER"`
	CertificateFile string        `yaml:"certificate_file" envconfig:"CERTIFICATE_FILE"`
	KeyFile         string        `yaml:"key_file" envconfig:"KEY_FILE"`
}

func Setup(file string) (*AppConfig, bool, error) {
	var created bool
	if _, err := os.Stat(file); os.IsNotExist(err) {
		content, err := assets.ReadConfigurationFile("app.dist.yaml")
		if err != nil {
			return nil, created, err
		}

		err = ioutil.WriteFile(file, content, 0744)
		if err != nil {
			return nil, created, err
		}

		created = true
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, created, fmt.Errorf("could not read configuration file: " + err.Error())
	}

	var conf AppConfig
	if err := yaml.Unmarshal(content, &conf); err != nil {
		return nil, created, fmt.Errorf("could not unmarshal configuration: " + err.Error())
	}

	if err := envconfig.Process("maestro", &conf); err != nil {
		return nil, created, fmt.Errorf("could not process env vars: " + err.Error())
	}

	return &conf, created, nil
}
