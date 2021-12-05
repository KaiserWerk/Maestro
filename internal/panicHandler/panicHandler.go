package panicHandler

import (
	"github.com/sirupsen/logrus"
)

func HandlePanic(l *logrus.Entry) {
	if r := recover(); r != nil {
		l.Panicf("recovered panic: %v\n", r)
	}
}
