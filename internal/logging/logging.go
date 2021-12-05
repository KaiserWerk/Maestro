package logging

import (
	"io"
	"os"
	"sync"

	"github.com/KaiserWerk/Maestro/internal/shutdownManager"

	"github.com/sirupsen/logrus"
)

type Mode uint8

const (
	ModeConsole Mode = iota
	ModeFile
	ModeBoth
)

var (
	err     error
	rotator *Rotator
)

func New(lvl logrus.Level, context string, mode Mode) *logrus.Entry {
	l := logrus.New()
	l.SetLevel(lvl)
	l.SetFormatter(&MaestroFormatter{
		LevelPadding:   7,
		ContextPadding: 9,
	})
	l.SetReportCaller(false)
	if mode == ModeConsole {
		l.SetOutput(os.Stdout)
	} else if mode == ModeFile {
		l.SetOutput(rotator)
	} else if mode == ModeBoth {
		l.SetOutput(io.MultiWriter(rotator, os.Stdout))
	}

	return l.WithField("context", context)
}

func Init(dir string) {
	shutdownManager.Register(CloseFileHandle)
	rotator, err = NewRotator(dir, "maestro.log", 10<<20, 0644)
	if err != nil {
		panic("cannot create rotator: " + err.Error())
	}
}

func CloseFileHandle(wg *sync.WaitGroup) {
	_ = rotator.Close()
	wg.Done()
}
