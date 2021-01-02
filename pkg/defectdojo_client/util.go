package defectdojo_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/sirupsen/logrus"
)

func BuildAuthHeader(api_key string) (string, string) {
	token_str := fmt.Sprintf("Token %s", api_key)
	return "Authorization", token_str
}

const APPLICATION_JSON = "application/json"

func (c *DefectdojoClient) BuildJsonRequestBytez(req_payload interface{}) (*bytes.Buffer, error) {
	bytez, err := json.Marshal(req_payload)
	if err != nil {
		return nil, fmt.Errorf("could not marshal to json: %s", err)
	}
	return bytes.NewBuffer(bytez), nil
}

func (c *DefectdojoClient) BuildMultipartFormBytez(values map[string]string, bytez []byte) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	f := multipart.NewWriter(&b)

	for k, v := range values {
		f.WriteField(k, v)
	}
	w, _ := f.CreatePart(textproto.MIMEHeader{
		"Content-Type":        []string{"text/xml"},
		"Content-Disposition": []string{`form-data; name="file"; filename="report.xml"`},
	})
	w.Write(bytez)
	f.Close()

	return &b, f.FormDataContentType(), nil
}

type PostErrors struct {
	NonFieldErrors []string `json:"non_field_errors"`
	FileErrors     []string `json:"file"`
}

func (c *DefectdojoClient) DoPost(api_path string, req_payload *bytes.Buffer, content_type string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v2/%s/", c.url, api_path)
	logrus.Debugf("POST %s", url)
	logrus.Debugf("data: %s", string(req_payload.Bytes()))

	req, err := http.NewRequest(http.MethodPost, url, req_payload)
	if err != nil {
		return nil, fmt.Errorf("something went wrong building request: %s", err)
	}
	req.Header.Set("Content-Type", content_type)

	logrus.Debugln("sending post")
	resp, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *DefectdojoClient) DoGet(api_path string, params map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v2/%s/", c.url, api_path)
	for k, v := range params {
		url = fmt.Sprintf("%s?%s=%s", url, k, v)
	}
	logrus.Debugf("GET %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("something went wrong building the request: %s", err)
	}

	return c.DoRequest(req)
}

func (c *DefectdojoClient) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add(BuildAuthHeader(c.api_key))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("recieved some kind of error: %s", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if resp.Body != nil {
			var e PostErrors
			decoder := json.NewDecoder(resp.Body)
			decoder.Decode(&e)
			logrus.Debugf("response that came back with error: %s, %s", e.NonFieldErrors, e.FileErrors)
		}
		return nil, fmt.Errorf("received status code of `%d`", resp.StatusCode)
	}

	return resp, nil
}
