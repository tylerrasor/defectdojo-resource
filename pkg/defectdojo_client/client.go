package defectdojo_client

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var application_json = "application/json"

type DefectdojoClient struct {
	client  *http.Client
	url     string
	api_key string
}

func NewDefectdojoClient(url string, api_key string) *DefectdojoClient {
	client := DefectdojoClient{
		client:  http.DefaultClient,
		url:     url,
		api_key: api_key,
	}

	return &client
}

func (c *DefectdojoClient) GetIdForProduct(name string) (int, error) {
	products, err := c.GetListOfProducts(name)
	if err != nil {
		return 0, err
	}

	// find the right product in the list
	logrus.Debugln(products)
	return products, nil
}

func (c *DefectdojoClient) GetListOfProducts(name string) (int, error) {
	url_path := fmt.Sprintf("%s/api/v2/products", c.url)
	req, err := http.NewRequest(http.MethodGet, url_path, nil)
	if err != nil {
		return 0, fmt.Errorf("something went wrong building request: %s", err)
	}

	resp, err := c.DoRequest(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %s", err)
	}
	logrus.Debugln(resp)

	return 1, nil
}

func (c *DefectdojoClient) GetOrCreateEngagement() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *DefectdojoClient) GetSomethingForIn() (string, error) {
	// I don't know what this will be, is this the 'most recent findings for report_type'?
	return "", fmt.Errorf("not implemented")
}
