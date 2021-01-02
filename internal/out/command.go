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
	engagement, err := c.CreateEngagement(p, request.Params.ReportType)
	if err != nil {
		return fmt.Errorf("error getting or creating engagement: %s", err)
	}
	logrus.Debugf("built new engagement, with id: %d", engagement.EngagementId)

	workdir := w.Args[1]
	full_path := fmt.Sprintf("%s/%s", workdir, request.Params.ReportPath)
	logrus.Debugf("trying to read file: %s", full_path)
	bytez, err := w.ReadFile(full_path)
	if err != nil {
		return fmt.Errorf("error reading report file: %s", err)
	}

	logrus.Debugln("uploading report")
	e, err := c.UploadReport(engagement.EngagementId, request.Params.ReportType, bytez)
	if err != nil {
		return fmt.Errorf("error uploading report: %s", err)
	}

	r := concourse.Response{
		Version: concourse.Version{
			Version: fmt.Sprint(e.EngagementId),
		},
	}
	if err := w.OutputResponseToConcourse(r); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
