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
		return nil, err
	}

	var results *ProductTypeSearchResults
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	if len(results.ProductTypeList) == 0 {
		return nil, fmt.Errorf("product `%s` not found", name)
	}
	if len(results.ProductTypeList) > 1 {
		return nil, fmt.Errorf("not sure how you did it, but got %d results for product_type name `%s`", len(results.ProductTypeList), name)
	}

	return &results.ProductTypeList[0], nil
}

func (c *DefectdojoClient) CreateProductType(name string) (*ProductType, error) {
	product_type_req := ProductType{
		Name: name,
	}

	payload, err := c.BuildJsonRequestBytez(product_type_req)
	if err != nil {
		return nil, err
	}
	resp, err := c.DoPost("product_types", payload, APPLICATION_JSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pt *ProductType
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&pt); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return pt, nil
}
