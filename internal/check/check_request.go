package check

import (
	"encoding/json"
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

type CheckRequest struct {
	Source  concourse.Source  `json:"source"`
	Version concourse.Version `json:"version"`
}

func (r CheckRequest) ValidateRequest() error {
	if err := r.Source.ValidateSource(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}

	return nil
}

func DecodeToCheckRequest(w *concourse.Worker) (*CheckRequest, error) {
	decoder := json.NewDecoder(w.Stdin)
	decoder.DisallowUnknownFields()

	var req CheckRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}

	if err := req.ValidateRequest(); err != nil {
		return nil, err
	}

	return &req, nil
}
