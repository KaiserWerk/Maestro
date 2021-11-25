package assets

import "embed"

//go:embed configuration
var configurationFS embed.FS

func ReadConfigurationFile(name string) ([]byte, error) {
	return configurationFS.ReadFile("configuration/" + name)
}
