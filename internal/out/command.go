package out

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Put(c *concourse.Concourse) error {
	logrus.SetOutput(c.Err)

	request, err := DecodeToPutRequest(c)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug logging on")
	}

	client := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)
	engagement_id, err := client.GetOrCreateEngagement()
	if err != nil {
		return fmt.Errorf("error getting or creating engagement: %s", err)
	}
	logrus.Debugln(engagement_id)

	// dump the response to stdout for concourse
	if err := concourse.OutputVersionToConcourse(c); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
