package concourse

import (
	"github.com/sirupsen/logrus"
)

func (w *Worker) setUpLogger() {
	l := logrus.New()
	l.SetOutput(w.stderr)
	l.SetLevel(logrus.ErrorLevel)
	l.SetFormatter(&logrus.TextFormatter{
		PadLevelText:     true,
		DisableTimestamp: true,
	})
	w.logger = l
}

func (w *Worker) EnableDebugLog() {
	w.logger.SetLevel(logrus.DebugLevel)
	w.logger.Debugln("debug logging on")
}

func (w *Worker) LogDebug(str string, args ...interface{}) {
	w.logger.Debugf(str, args...)
}

func (w *Worker) LogError(str string, args ...interface{}) {
	w.logger.Errorf(str, args...)
}
