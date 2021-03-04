package in

import (
	"fmt"
)

type GetParams struct {
	ReportType string `json:"report_type"`
}

func (p GetParams) ValidateParams() error {
	if p.ReportType == "" {
		return fmt.Errorf("Required parameter `report_type` not supplied.")
	}
	return nil
}
