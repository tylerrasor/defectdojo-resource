package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

var application_json = "application/json"

type DefectdojoClient struct {
	url              string
	api_key          string
	usernamePassword userPass
}

type userPass struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDefectdojoClient(url string, username string, password string, api_key string) (*DefectdojoClient, error) {
	dd_client := DefectdojoClient{
		url:     url,
		api_key: api_key,
		usernamePassword: userPass{
			Username: username,
			Password: password,
		},
	}

	logrus.Debugln("trying to authenticate to defectdojo instance")
	if err := dd_client.tryToAuth(); err != nil {
		return nil, fmt.Errorf("error in authenticating to defectdojo instance at `%s`: %s", dd_client.url, err)
	}

	return &dd_client, nil
}

func (c *DefectdojoClient) tryToAuth() error {
	if c.api_key == "" {
		logrus.Debugln("no api_key provided, trying to auth with user/pass instead")
		if err := c.getApiKey(); err != nil {
			return fmt.Errorf("something went wrong trying to get api_key for user `%s`: %s", c.usernamePassword.Username, err)
		}
	}

	url_path := fmt.Sprintf("%s/api/v2/products", c.url)
	req, err := http.NewRequest(http.MethodGet, url_path, nil)
	if err != nil {
		return fmt.Errorf("something went wrong building request: %s", err)
	}

	token_str := fmt.Sprintf("Token %s", c.api_key)
	req.Header.Add("Authorization", token_str)

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

// todo: test this and fix it
func (c *DefectdojoClient) getApiKey() error {
	url_path := fmt.Sprintf("%s/api/v2/api-token-auth", c.url)
	payload, err := json.Marshal(c.usernamePassword)
	if err != nil {
		return fmt.Errorf("error marshalling to JSON: %s", err)
	}

	resp, err := http.Post(url_path, application_json, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("recieved some kind of error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received status code of `%d`", resp.StatusCode)
	}

	// not 100% sure how to decode this return type, skipping for now
	// { token: 'api_key' }
	var body []byte
	resp.Body.Read(body)

	return nil
}

func (c *DefectdojoClient) GetOrCreateEngagement() (int, error) {
	return 0, fmt.Errorf("not implemented")
}
