package in

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/client"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
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

	logrus.Debugln("creating http client")
	client, err := client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.Username, request.Source.Password, request.Source.ApiKey)
	if err != nil {
		return fmt.Errorf("error creating client to interact with defectdojo: %s", err)
	}

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
