package check

import (
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func Check(w *concourse.Worker) error {
	request, err := DecodeToCheckRequest(w)
	if err != nil {
		return fmt.Errorf("invalid payload: %s", err)
	}

	if request.Source.Debug {
		w.EnableDebugLog()
	}

	client := defectdojo_client.NewDefectdojoClient(request.Source.DefectDojoUrl, request.Source.ApiKey)

	issue_url := "https://github.com/tylerrasor/defectdojo-resource/issues/29"
	w.LogDebug("trying to find engagement_id `%s`", request.Version.EngagementId)
	var r []concourse.Version
	if request.Version.EngagementId == "" || request.Version.EngagementId == issue_url {
		w.LogDebug("need to implement `go get latest engagement of report_type`")
		r = []concourse.Version{
			{
				EngagementId: issue_url,
			},
		}
	} else {
		e, err := client.GetEngagement(request.Version.EngagementId)
		if err != nil {
			return fmt.Errorf("error getting engagement: %s", err)
		}

		w.LogDebug("building response object")
		r = []concourse.Version{
			{
				EngagementId: fmt.Sprint(e.EngagementId),
			},
		}
	}

	if err := w.OutputResponseToConcourse(r); err != nil {
		return fmt.Errorf("error encoding response to JSON: %s", err)
	}

	return nil
}
