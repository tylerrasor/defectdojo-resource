package concourse

import (
	"fmt"
	"strings"
)

type Response struct {
	Version Version `json:"version"`
}

type Version struct {
	Version string `json:"version"`
}

type Source struct {
	DefectDojoUrl string `json:"defectdojo_url"`
	ApiKey        string `json:"api_key"`
	Debug         bool   `json:"debug"`
}

func (s *Source) ValidateSource() error {
	if s.DefectDojoUrl == "" {
		return fmt.Errorf("Required `defectdojo_url` not supplied.")
	}
	if !strings.HasPrefix(s.DefectDojoUrl, "http://") && !strings.HasPrefix(s.DefectDojoUrl, "https://") {
		return fmt.Errorf("Please provide http(s):// prefix in `defectdojo_url`.")
	}
	if s.ApiKey == "" {
		return fmt.Errorf("Required `api_key` not supplied.")
	}
	return nil
}
