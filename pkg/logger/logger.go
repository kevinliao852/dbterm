package logger

import (
	log "github.com/sirupsen/logrus"
	"io"
)

type LoggerOption struct {
	log    *log.Logger
	prefix string
}

func NewLoggerOption(log *log.Logger) *LoggerOption {
	return &LoggerOption{
		log:    log,
		prefix: "",
	}
}

func (l *LoggerOption) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func (l *LoggerOption) SetPrefix(s string) {
	l.prefix = s
}
