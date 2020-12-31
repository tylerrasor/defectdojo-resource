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

func (c *DefectdojoClient) GetIdForProduct(name string) error {
	url_path := fmt.Sprintf("%s/api/v2/products", c.url)
	req, err := http.NewRequest(http.MethodGet, url_path, nil)
	if err != nil {
		return fmt.Errorf("something went wrong building request: %s", err)
	}

	token_str := fmt.Sprintf("Token %s", c.api_key)
	req.Header.Add("Authorization", token_str)

	logrus.Debugln("grabbing list of products just to confirm we can auth")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("recieved some kind of error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code of `%d`", resp.StatusCode)
	}

	return nil
}

func (c *DefectdojoClient) GetOrCreateEngagement() (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (c *DefectdojoClient) GetSomethingForIn() (string, error) {
	// I don't know what this will be, is this the 'most recent findings for report_type'?
	return "", fmt.Errorf("not implemented")
}
