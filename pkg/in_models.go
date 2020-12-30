package resource

import (
	"fmt"
)

type GetRequest struct {
	Source Source    `json:"source"`
	Params GetParams `json:"params"`
}

func (g GetRequest) Validate() error {
	return fmt.Errorf("not implemented yet")
}

type GetParams struct {
}

func (p GetParams) ValidateParams() error {
	return fmt.Errorf("not implemented yet")
}
