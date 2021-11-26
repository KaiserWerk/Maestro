package handler

import (
	"github.com/sirupsen/logrus"
)

type HttpHandler struct {
	Logger *logrus.Entry
}

func (h HttpHandler) ContextLogger(context string) *logrus.Entry {
	return h.Logger.WithField("context", context)
}
