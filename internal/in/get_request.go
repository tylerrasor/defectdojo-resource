package in

import (
	"encoding/json"
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

type GetRequest struct {
	Source concourse.Source `json:"source"`
	Params GetParams        `json:"params"`
}

func (r GetRequest) ValidateRequest() error {
	if err := r.Source.ValidateSource(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}

	if err := r.Params.ValidateParams(); err != nil {
		return fmt.Errorf("invalid params config: %s", err)
	}

	return nil
}

func DecodeToGetRequest(w *concourse.Worker) (*GetRequest, error) {
	decoder := json.NewDecoder(w.Stdin)
	decoder.DisallowUnknownFields()

	var req GetRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}

	if err := req.ValidateRequest(); err != nil {
		return nil, err
	}

	return &req, nil
}
