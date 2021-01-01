package out

import (
	"fmt"

	"github.com/tylerrasor/defectdojo-resource/internal/util"
)

type PutParams struct {
	ReportType string `json:"report_type"`
	ReportPath string `json:"path_to_report"`
}

func (p PutParams) ValidateParams() error {
	if p.ReportType == "" {
		return fmt.Errorf("Required parameter `report_type` not supplied.")
	}

	implemented, ok := util.SupportedReportTypes[p.ReportType]
	if !ok {
		return fmt.Errorf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", p.ReportType)
	}
	if !implemented {
		return fmt.Errorf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", p.ReportType)
	}

	if p.ReportPath == "" {
		return fmt.Errorf("Required parameter `path_to_report` not supplied.")
	}

	return nil
}
