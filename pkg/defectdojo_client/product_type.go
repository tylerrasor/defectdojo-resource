package defectdojo_client

import (
	"encoding/json"
	"fmt"
)

type ProductTypeSearchResults struct {
	ProductTypeList []ProductType `json:"results"`
}

type ProductType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (c *DefectdojoClient) GetProductType(name string) (*ProductType, error) {
	params := map[string]string{
		"name": name,
	}
	resp, err := c.DoGet("product_types", params)
	if err != nil {
		return false, err
	}

	var results *ProductTypeSearchResults
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		return nil, err
	}

	if len(results.ProductTypeList) == 0 {
		return nil, fmt.Errorf("product `%s` not found", name)
	}
	if len(results.ProductTypeList) > 1 {
		return nil, fmt.Errorf("not sure how you did it, but got %d results for product_type name `%s`", len(results.ProductTypeList), name)
	}

	return &results.ProductTypeList[0], nil
}
