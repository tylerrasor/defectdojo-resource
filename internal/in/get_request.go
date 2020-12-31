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

func (g GetRequest) Validate() error {
	return fmt.Errorf("not implemented yet")
}

func DecodeToGetRequest(w *concourse.Worker) (*GetRequest, error) {
	decoder := json.NewDecoder(w.In)
	decoder.DisallowUnknownFields()

	var req GetRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
