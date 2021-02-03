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
	DefectDojoUrl            string `json:"defectdojo_url"`
	ApiKey                   string `json:"api_key"`
	ProductName              string `json:"product_name"`
	CreateProductIfNotExists bool   `json:"create_product_if_not_exist"`
	ProductType              string `json:"product_type"`
	Debug                    bool   `json:"debug"`
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
	if s.ProductName == "" {
		return fmt.Errorf("Required `product_name` not supplied.")
	}
	if s.CreateProductIfNotExists {
		if s.ProductType == "" {
			return fmt.Errorf("Optional `create_product_if_not_exist` set and required `product_type` not supplied.")
		}
	}
	return nil
}
