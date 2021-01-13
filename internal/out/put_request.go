package out

import (
	"encoding/json"
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

type PutRequest struct {
	Source concourse.Source `json:"source"`
	Params PutParams        `json:"params"`
}

func (r PutRequest) ValidateRequest() error {
	if err := r.Source.ValidateSource(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}

	if err := r.Params.ValidateParams(); err != nil {
		return fmt.Errorf("invalid params config: %s", err)
	}

	return nil
}

func DecodeToPutRequest(w *concourse.Worker) (*PutRequest, error) {
	decoder := json.NewDecoder(w.Stdin)
	decoder.DisallowUnknownFields()

	w.LogDebug("decoding the concourse input to put request")
	var req PutRequest
	if err := decoder.Decode(&req); err != nil {
		return nil, err
	}

	if err := req.ValidateRequest(); err != nil {
		return nil, err
	}

	return &req, nil
}
