package middleware

import (
	"github.com/sirupsen/logrus"

	"github.com/KaiserWerk/Maestro/internal/configuration"
)

type MWHandler struct {
	Config *configuration.AppConfig
	Logger *logrus.Entry
}
