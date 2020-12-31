package in

import (
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

type GetParams struct {
}

func (p GetParams) ValidateParams() error {
	return fmt.Errorf("not implemented yet")
}
