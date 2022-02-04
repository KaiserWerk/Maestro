package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
)

type Mode uint8

const (
	ModeConsole Mode = iota
	ModeFile
	ModeBoth
)

const (
	logFile             = "maestro.log"
	maxSize             = 10 << 20
	perms   fs.FileMode = 0644
)

func New(dir, context string, lvl logrus.Level, mode Mode) (*logrus.Entry, func() error, error) {
	var (
		l  = logrus.New()
		cf func() error
	)
	l.SetLevel(lvl)
	l.SetFormatter(&MaestroFormatter{
		LevelPadding:   7,
		ContextPadding: 9,
	})
	l.SetReportCaller(false)
	if mode == ModeConsole {
		l.SetOutput(os.Stdout)
	} else if mode == ModeFile {
		rotator, err := NewRotator(dir, logFile, maxSize, perms)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create rotator: " + err.Error())
		}
		cf = func() error {
			return rotator.Close()
		}
		l.SetOutput(rotator)
	} else if mode == ModeBoth {
		rotator, err := NewRotator(dir, logFile, maxSize, perms)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot create rotator: " + err.Error())
		}
		cf = func() error {
			return rotator.Close()
		}
		l.SetOutput(io.MultiWriter(rotator, os.Stdout))
	}

	return l.WithField("context", context), cf, nil
}
