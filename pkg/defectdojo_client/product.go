package defectdojo_client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ProductSearchResults struct {
	ProductList []Product `json:"results"`
}

type Product struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

func (c *DefectdojoClient) GetProduct(name string) (*Product, error) {
	// get list of products
	url := fmt.Sprintf("%s/api/v2/products", c.url)
	logrus.Debugf("GET %s\n", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("something went wrong building request: %s", err)
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var results *ProductSearchResults
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	logrus.Debugf("found %d products\n", len(results.ProductList))
	var p *Product
	for i := range results.ProductList {
		logrus.Debugf("product name: %s\n", results.ProductList[i].Name)
		if results.ProductList[i].Name == name {
			logrus.Debugln("found product in list")
			p = &results.ProductList[i]
			break
		}
	}

	if p == nil {
		return nil, fmt.Errorf("product `%s` not found", name)
	}

	return p, nil
}

func (c *DefectdojoClient) CreateEngagement(p *Product) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
