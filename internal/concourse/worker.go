package concourse

import (
	"encoding/json"
	"io"
)

func AttachToWorker(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Worker {
	return &Worker{
		In:   stdin,
		Err:  stderr,
		Out:  stdout,
		Args: args,
	}
}

type Worker struct {
	In   io.Reader
	Err  io.Writer
	Out  io.Writer
	Args []string
}

func (w *Worker) OutputResponseToConcourse(r Response) error {
	return json.NewEncoder(w.Out).Encode(r)
}
