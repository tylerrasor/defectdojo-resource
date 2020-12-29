package resource

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/client"
	"github.com/tylerrasor/defectdojo-resource/models"
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

func (o *Out) Execute() error {
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

func DecodeFromOut(o *Out) (*models.PutRequest, error) {
	var request models.PutRequest

	decoder := json.NewDecoder(o.stdin)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return nil, err
	}
	return &request, nil
}

func OutputVersionToConcourse(o *Out) error {
	version_str := "need to figure out unique combination of app name, version, build number, something"
	message := fmt.Sprintf("preparing to JSON encode response: %s", version_str)
	logrus.Debugln(message)

	version := models.Version{
		Version: version_str,
	}
	response := models.PutResponse{
		Version: version,
	}

	return json.NewEncoder(o.stdout).Encode(response)
}
