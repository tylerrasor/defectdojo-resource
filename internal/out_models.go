package resource

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type PutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

func (r PutRequest) Validate() error {
	logrus.Debugln("getting ready to validate source")
	if err := r.Source.ValidateSource(); err != nil {
		return fmt.Errorf("invalid source config: %s", err)
	}

	logrus.Debugln("getting ready to validate params")
	if err := r.Params.ValidateParams(); err != nil {
		return fmt.Errorf("invalid params config: %s", err)
	}

	return nil
}

type PutParams struct {
	ReportType string `json:"report_type"`
}

func (p PutParams) ValidateParams() error {
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
