package models

import (
	"fmt"
	"strings"
)

type Source struct {
	DefectDojoUrl string `json:"defectdojo_url"`
	Username      string `json:"username"`
	Password      string `json:"password,omitempty"`
	ApiKey        string `json:"api_key"`
	Debug         bool   `json:"debug"`
}

func (s *Source) Validate() error {
	if s.DefectDojoUrl == "" {
		return fmt.Errorf("Required `defectdojo_url` not supplied.")
	}
	if !strings.HasPrefix(s.DefectDojoUrl, "http://") && !strings.HasPrefix(s.DefectDojoUrl, "https://") {
		return fmt.Errorf("Please provide http(s):// prefix in `defectdojo_url`.")
	}
	if s.Username == "" {
		return fmt.Errorf("Required `username` not supplied.")
	}
	if s.ApiKey == "" {
		return fmt.Errorf("Required `api_key` not supplied.")
	}
	return nil
}

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

type Version struct {
	Version string `json:"version"`
}

type PutResponse struct {
	Version Version `json:"version"`
}
