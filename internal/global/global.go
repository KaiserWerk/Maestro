package global

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"sync/atomic"
)

const (
	DefaultPort  = "9200"
	AuthHeader   = "X-Registry-Token"
	EnvAuthToken = "MAESTRO_TOKEN"
)

var (
	authToken = ""
	testPort uint32 = 30000
)

func SetAuthToken(t string) {
	if t != "" {
		authToken = t
	} else if os.Getenv(EnvAuthToken) != "" {
		authToken = os.Getenv(EnvAuthToken)
	}

	if authToken == "" {
		panic("authToken is empty!")
	}
}

func GetAuthToken() string {
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

func GetPortForTest() uint32 {
	return atomic.AddUint32(&testPort, 1)
}