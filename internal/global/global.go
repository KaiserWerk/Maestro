package global

import (
	"crypto/rand"
	"encoding/base64"
	"sync/atomic"
)

const (
	DefaultPort = "9200"
	AuthHeader  = "X-Registry-Token"

	envPrefix      = "MAESTRO_"
	EnvBindAddress = envPrefix + "BIND_ADDRESS"
	EnvAuthToken   = envPrefix + "TOKEN"
	EnvDieAfter    = envPrefix + "DIE_AFTER"
	EnvCertFile    = envPrefix + "CERT_FILE"
	EnvKeyFile     = envPrefix + "KEY_FILE"
)

var (
	testPort uint32 = 30000
)

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
