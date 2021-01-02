package defectdojo_client

import (
	"encoding/json"
	"fmt"
)

type ProductSearchResults struct {
	ProductList []Product `json:"results"`
}

type Product struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (c *DefectdojoClient) GetProduct(name string) (*Product, error) {
	params := map[string]string{
		"name": name,
	}
	resp, err := c.DoGet("products", params)
	if err != nil {
		return nil, err
	}

	var results *ProductSearchResults
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	if len(results.ProductList) == 0 {
		return nil, fmt.Errorf("product `%s` not found", name)
	}
	if len(results.ProductList) > 1 {
		return nil, fmt.Errorf("not sure how you did it, but got %d results for product name `%s`", len(results.ProductList), name)
	}

	return &results.ProductList[0], nil
}
