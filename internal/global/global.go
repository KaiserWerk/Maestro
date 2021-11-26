package global

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

const (
	DefaultPort = "9200"
	AuthHeader  = "X-Registry-Token"
	PingHeader  = "X-Ping-Token"

	EnvAuthToken = "MAESTRO_TOKEN"
)

var (
	authToken = ""
)

func SetToken(t string) {
	if t != "" {
		authToken = t
	} else if os.Getenv(EnvAuthToken) != "" {
		authToken = os.Getenv(EnvAuthToken)
	}
}

func GetToken() string {
	return authToken
}

func GenerateToken(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(b)
}