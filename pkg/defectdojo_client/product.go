package defectdojo_client

import (
	"bytes"
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
	url := fmt.Sprintf("%s/api/v2/products/?name=%s", c.url, name)
	logrus.Debugf("GET %s", url)
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

	if len(results.ProductList) == 0 {
		return nil, fmt.Errorf("product `%s` not found", name)
	}
	if len(results.ProductList) > 1 {
		return nil, fmt.Errorf("not sure how you did it, but got %d results for product name `%s`", len(results.ProductList), name)
	}

	return &results.ProductList[0], nil
}

type Engagement struct {
	EngagementId   int    `json:"id,omitempty"`
	ProductId      int    `json:"product"`
	StartDate      string `json:"target_start"`
	EndDate        string `json:"target_end"`
	EngagementType string `json:"engagement_type"`
}

func (c *DefectdojoClient) CreateEngagement(p *Product) (*Engagement, error) {
	url := fmt.Sprintf("%s/api/v2/engagements/", c.url)
	logrus.Debugf("POST %s", url)

	engagement_req := Engagement{
		ProductId:      p.Id,
		StartDate:      "2021-01-01",
		EndDate:        "2021-01-01",
		EngagementType: "CI/CD",
	}
	bytez, err := json.Marshal(engagement_req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal to json: %s", err)
	}

	logrus.Debugf("trying to send payload: %s", string(bytez))
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytez))
	if err != nil {
		return nil, fmt.Errorf("something went wrong building request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")

	logrus.Debugln("sending post")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var e *Engagement
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&e); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err)
	}

	return e, nil
}
