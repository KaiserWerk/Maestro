package entity

type AppConfig struct {
	BindAddress     string `yaml:"bind_address" envconfig:"BIND_ADDRESS"`
	AuthToken       string `yaml:"auth_token" envconfig:"AUTH_TOKEN"`
	DieAfter        uint   `yaml:"die_after" envconfig:"DIE_AFTER"`
	CertificateFile string `yaml:"certificate_file" envconfig:"CERTIFICATE_FILE"`
	KeyFile         string `yaml:"key_file" envconfig:"KEY_FILE"`
}
