package in

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Get(c *concourse.Concourse) error {
	logrus.SetOutput(c.Err)

	request, err := DecodeToGetRequest(c)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug logging on")
	}

	client := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)
	something, err := client.GetSomethingForIn()
	if err != nil {
		return fmt.Errorf("error getting something: %s", err)
	}
	logrus.Debugln(something)

	if err := concourse.OutputVersionToConcourse(c); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
