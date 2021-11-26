package entity

type AppConfig struct {
	BindAddress     string `yaml:"bind_address"`
	CertificateFile string `yaml:"certificate_file"`
	KeyFile         string `yaml:"key_file"`
}
