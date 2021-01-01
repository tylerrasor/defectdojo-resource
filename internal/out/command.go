package out

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Put(w *concourse.Worker) error {
	logrus.SetOutput(w.Err)

	request, err := DecodeToPutRequest(w)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("debug logging on")
	}

	c := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)

	logrus.Debugln("looking for product profile")
	p, err := c.GetProduct(request.Source.AppName)
	if err != nil {
		return fmt.Errorf("error getting product: %s", err)
	}

	logrus.Debugln("creating new cicd engagement")
	engagement_id, err := c.CreateEngagement(p)
	if err != nil {
		return fmt.Errorf("error getting or creating engagement: %s", err)
	}
	logrus.Debugln(engagement_id)

	logrus.Debugln("building response")
	r := concourse.Response{
		Version: concourse.Version{
			Version: "need to figure out unique combination of app name, version, build number, something",
		},
	}
	if err := w.OutputResponseToConcourse(r); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
