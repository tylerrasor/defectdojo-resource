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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("something went wrong building request: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recieved some kind of error: %s", err)
	}
	defer resp.Body.Close()

	var results *ProductSearchResults
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	var p *Product
	for i := range results.ProductList {
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

func (c *DefectdojoClient) GetOrCreateEngagement(p *Product) (int, error) {
	return 0, fmt.Errorf("not implemented")
}
