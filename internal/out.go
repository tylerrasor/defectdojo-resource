package resource

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/client"
)

func (c *Concourse) Put() error {
	logrus.SetOutput(c.stderr)

	request, err := DecodeToPutRequest(c)
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

	engagement_id, err := client.GetOrCreateEngagement()
	if err != nil {
		return fmt.Errorf("error getting or creating engagement: %s", err)
	}
	logrus.Debugln(engagement_id)

	// dump the response to stdout for concourse
	if err := OutputVersionToConcourse(c); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}

func DecodeToPutRequest(c *Concourse) (*PutRequest, error) {
	decoder := json.NewDecoder(c.stdin)
	decoder.DisallowUnknownFields()

	var req PutRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
