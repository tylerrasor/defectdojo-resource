package check

import (
	"encoding/json"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

type CheckRequest struct {
	Source  concourse.Source  `json:"source"`
	Version concourse.Version `json:"version"`
}

func DecodeToCheckRequest(w *concourse.Worker) (*CheckRequest, error) {
	decoder := json.NewDecoder(w.Stdin)
	decoder.DisallowUnknownFields()

	var req CheckRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
