package defectdojo_client

import (
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
