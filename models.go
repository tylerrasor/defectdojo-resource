package resource

import (
	"errors"
	"strings"
)

type Source struct {
	DefectDojoUrl string `json:"defectdojo_url"`
	Debug         bool   `json:"debug"`
}

func (s *Source) Validate() error {
	if s.DefectDojoUrl == "" {
		return errors.New("Required parameter `DefectDojoUrl` not supplied.")
	}
	if !strings.HasPrefix(s.DefectDojoUrl, "http://") || !strings.HasPrefix(s.DefectDojoUrl, "https://") {
		return errors.New("Please provide http(s):// prefix")
	}
	return nil
}
