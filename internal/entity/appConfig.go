package entity

type AppConfig struct {
	App struct {
		BindAddress     string `yaml:"bind_address"`
		AuthToken       string `yaml:"auth_token"`
		DieAfter        int    `yaml:"die_after"`
		CertificateFile string `yaml:"certificate_file"`
		KeyFile         string `yaml:"key_file"`
	} `yaml:"app"`
}
