package defectdojo_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func BuildAuthHeader(api_key string) (string, string) {
	token_str := fmt.Sprintf("Token %s", api_key)
	return "Authorization", token_str
}

// func (c *DefectdojoClient) DoGet(resp_obj Response, api_path string, query_params ...string) (*Response, error) {
// 	url := fmt.Sprintf("%s/api/v2/%s/", c.url, api_path)
// 	for i := range query_params {
// 		if i%2 == 0 {
// 			url = fmt.Sprintf("?%s=%s", url, query_params[i], query_params[i+1])
// 		}
// 	}
// 	logrus.Debugf("GET %s", url)

// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("something went wrong building request: %s", err)
// 	}

// 	resp, err := c.DoRequest(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	results := resp_obj.NewOfType()
// 	decoder := json.NewDecoder(resp.Body)
// 	if err := decoder.Decode(&results); err != nil {
// 		return nil, fmt.Errorf("error decoding response: %s", err)
// 	}

// 	return results, nil
// }

func (c *DefectdojoClient) DoPost(api_path string, req_payload interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v2/%s/", c.url, api_path)
	logrus.Debugf("POST %s", url)

	bytez, err := json.Marshal(req_payload)
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

	return resp, nil
}

func (c *DefectdojoClient) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add(BuildAuthHeader(c.api_key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recieved some kind of error: %s", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received status code of `%d`", resp.StatusCode)
	}

	return resp, nil
}
