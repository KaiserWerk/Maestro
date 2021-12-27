package handler

import (
	"github.com/KaiserWerk/Maestro/internal/cache"

	"github.com/sirupsen/logrus"
)

type HttpHandler struct {
	Logger       *logrus.Entry
	MaestroCache *cache.MaestroCache
}

func (h HttpHandler) ContextLogger(context string) *logrus.Entry {
	return h.Logger.WithField("context", context)
}
