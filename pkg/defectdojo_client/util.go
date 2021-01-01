package defectdojo_client

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func BuildAuthHeader(api_key string) (string, string) {
	token_str := fmt.Sprintf("Token %s", api_key)
	return "Authorization", token_str
}

func (c *DefectdojoClient) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add(BuildAuthHeader(c.api_key))

	logrus.Debugln("grabbing list of products just to confirm we can auth")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recieved some kind of error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code of `%d`", resp.StatusCode)
	}

	return &http.Response{}, nil
}
