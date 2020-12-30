package resource

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/client"
)

func (o *Concourse) Put() error {
	logrus.SetOutput(o.stderr)

	request, err := DecodeFromOut(o)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug logging on")
	}

	logrus.Debugln("getting ready to validate source")
	if err := request.Source.Validate(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}

	logrus.Debugln("getting ready to validate params")
	if err := request.Params.Validate(); err != nil {
		return fmt.Errorf("invalid params config: %s", err)
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
	if err := OutputVersionToConcourse(o); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}

func DecodeFromOut(o *Concourse) (*PutRequest, error) {
	var request PutRequest

	decoder := json.NewDecoder(o.stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}
