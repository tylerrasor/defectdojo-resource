package out

import (
	"fmt"
)

type PutParams struct {
	ReportType      string `json:"report_type"`
	ReportPath      string `json:"path_to_report"`
	CloseEngagement bool   `json:"close_engagement"`
}

func (p PutParams) ValidateParams() error {
	if p.ReportType == "" {
		return fmt.Errorf("Required parameter `report_type` not supplied.")
	}

	if p.ReportPath == "" {
		return fmt.Errorf("Required parameter `path_to_report` not supplied.")
	}

	return nil
}
