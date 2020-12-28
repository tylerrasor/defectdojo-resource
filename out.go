package resource

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
)

type Out struct {
	stdin  io.Reader
	stderr io.Writer
	stdout io.Writer
	args   []string
}

func NewOut(
	stdin io.Reader,
	stderr io.Writer,
	stdout io.Writer,
	args []string,
) *Out {
	return &Out{
		stdin:  stdin,
		stderr: stderr,
		stdout: stdout,
		args:   args,
	}
}

type PutParams struct {
}

type PutRequest struct {
	Source    Source
	PutParams PutParams
}

func (o *Out) Execute() error {
	var request PutRequest

	decoder := json.NewDecoder(o.stdin)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&request)

	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug loggin on")
	}

	return nil
}
