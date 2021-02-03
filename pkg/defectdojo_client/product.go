package defectdojo_client

import (
	"encoding/json"
	"fmt"
)

type ProductSearchResults struct {
	ProductList []Product `json:"results"`
}

type Product struct {
	Name          string `json:"name"`
	Id            int    `json:"id,omitempty"`
	ProductTypeId int    `json:"prod_type"`
	Description   string `json:"description"`
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

func (c *DefectdojoClient) CreateProduct(name string, product_type string) (*Product, error) {
	pt, err := c.GetProductType(product_type)
	if err != nil {
		return nil, err
	}

	product_req := Product{
		Name:          name,
		ProductTypeId: pt.Id,
		Description:   "created by CI/CD",
	}

	payload, err := c.BuildJsonRequestBytez(product_req)
	if err != nil {
		return nil, err
	}
	resp, err := c.DoPost("products", payload, APPLICATION_JSON)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var p *Product
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&p); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return p, nil
}
