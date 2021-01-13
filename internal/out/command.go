package out

import (
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Put(w *concourse.Worker) error {
	request, err := DecodeToPutRequest(w)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		w.EnableDebugLog()
	}

	c := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)

	w.LogDebug("looking for product profile")
	p, err := c.GetProduct(request.Source.AppName)
	if err != nil {
		return fmt.Errorf("error getting product: %s", err)
	}

	w.LogDebug("creating new cicd engagement")
	engagement, err := c.CreateEngagement(p, request.Params.ReportType)
	if err != nil {
		return fmt.Errorf("error getting or creating engagement: %s", err)
	}
	w.LogDebug("built new engagement, with id: %d", engagement.EngagementId)

	workdir := w.GetWorkDir()
	full_path := fmt.Sprintf("%s/%s", workdir, request.Params.ReportPath)
	w.LogDebug("trying to read file: %s", full_path)
	bytez, err := w.ReadFile(full_path)
	if err != nil {
		return fmt.Errorf("error reading report file: %s", err)
	}

	w.LogDebug("uploading report")
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
