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

	client := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)
	engagement_id, err := client.GetOrCreateEngagement()
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
