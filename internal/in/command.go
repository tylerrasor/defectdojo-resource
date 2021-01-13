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
	something, err := client.GetSomethingForIn()
	if err != nil {
		return fmt.Errorf("error getting something: %s", err)
	}
	w.LogDebug(something)

	w.LogDebug("building response object")
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
