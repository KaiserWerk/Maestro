package entity

type AppConfig struct {
	App struct {
		BindAddress     string `yaml:"bind_address"`
		AuthToken       string `yaml:"auth_token"`
		CertificateFile string `yaml:"certificate_file"`
		KeyFile         string `yaml:"key_file"`
	} `yaml:"app"`
	Database struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"database"`
}
