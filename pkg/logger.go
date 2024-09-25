package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

const p string = "./logs.json"

var Log *logrus.Logger = initLogrus(p)

type hook struct {
	Writers []io.Writer
	HLevels []logrus.Level
}

func (h *hook) Fire(e *logrus.Entry) error {

	bytes, err := e.Bytes()

	if err != nil {
		return err
	}

	for _, h := range h.Writers {
		if _, err := h.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

func (h *hook) Levels() []logrus.Level {
	return h.HLevels
}

func initLogrus(logFilePath string) *logrus.Logger {

	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s: %d", path.Base(f.File), f.Line)
		}})

	h := &hook{
		Writers: []io.Writer{},
		HLevels: logrus.AllLevels,
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Panic(err)
		return nil
	}

	l.SetOutput(io.Discard)

	h.Writers = append(h.Writers, file)

	l.AddHook(h)
	l.SetLevel(logrus.TraceLevel)

	return l

}
