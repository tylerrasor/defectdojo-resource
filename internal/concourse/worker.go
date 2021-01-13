package concourse

import (
	"encoding/json"
	"io"

	"github.com/sirupsen/logrus"
)

func AttachToWorker(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Worker {
	w := Worker{
		Stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
	w.setUpLogger()
	return &w
}

type Worker struct {
	Stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
	logger *logrus.Logger
}

func (w *Worker) OutputResponseToConcourse(r Response) error {
	return json.NewEncoder(w.stdout).Encode(r)
}

func (w *Worker) GetWorkDir() string {
	return w.args[1]
}
