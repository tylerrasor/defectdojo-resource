package in

import (
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Get(w *concourse.Worker) error {
	request, err := DecodeToGetRequest(w)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		w.EnableDebugLog()
	}

	client := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)
	p, err := client.GetProduct(request.Source.ProductName)
	if err != nil {
		return fmt.Errorf("error getting product: %s", err)
	}
	w.LogDebug("found product, with id: %d", p.Id)

	e, err := client.GetEngagementForReportType(p, request.Params.ReportType)
	if err != nil {
		return fmt.Errorf("error getting engagement: %s", err)
	}
	w.LogDebug("found engagement, with id: %d", e.EngagementId)

	w.LogDebug("building response object")
	r := concourse.Response{
		Version: concourse.Version{
			EngagementId: fmt.Sprint(e.EngagementId),
		},
	}
	if err := w.OutputResponseToConcourse(r); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
