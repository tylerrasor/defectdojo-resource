package defectdojo_client

import (
	"fmt"
	"net/http"
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

func (c *DefectdojoClient) GetSomethingForIn() (string, error) {
	// I don't know what this will be, is this the 'most recent findings for report_type'?
	return "", fmt.Errorf("not implemented")
}
