package resource

import "fmt"

/* specific to PUT */
type PutParams struct {
	ReportType string `json:"report_type"`
}

func (p *PutParams) Validate() error {
	if p.ReportType == "" {
		return fmt.Errorf("Required parameter `report_type` not supplied.")
	}

	implemented, key_exists := SupportedReportTypes[p.ReportType]
	if !key_exists {
		return fmt.Errorf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", p.ReportType)
	}
	if !implemented {
		return fmt.Errorf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", p.ReportType)
	}

	return nil
}

type PutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}
