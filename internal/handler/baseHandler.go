package handler

import (
	"github.com/KaiserWerk/Maestro/internal/cache"
	"github.com/KaiserWerk/Maestro/internal/configuration"

	"github.com/sirupsen/logrus"
)

type BaseHandler struct {
	Config       *configuration.AppConfig
	Logger       *logrus.Entry
	MaestroCache *cache.MaestroCache
}

func (h BaseHandler) ContextLogger(context string) *logrus.Entry {
	return h.Logger.WithField("context", context)
}
