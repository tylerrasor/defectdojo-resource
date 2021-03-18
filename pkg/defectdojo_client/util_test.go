package defectdojo_client_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func TestBuildAuthHeader(t *testing.T) {
	key := "api_key"

	k, v := defectdojo_client.BuildAuthHeader(key)

	assert.Equal(t, k, "Authorization")
	token_str := fmt.Sprintf("Token %s", key)
	assert.Equal(t, v, token_str)
}

func TestBuildJsonRequestBytez(t *testing.T) {
	id := 5
	target_date := "2021-01-01"
	e_type := "type"
	name := "engagement name"
	payload := defectdojo_client.Engagement{
		ProductId:      id,
		StartDate:      target_date,
		EndDate:        target_date,
		EngagementType: e_type,
		EngagementName: name,
	}

	c := defectdojo_client.NewDefectdojoClient("nil", "nil")
	bytez, err := c.BuildJsonRequestBytez(payload)

	json := fmt.Sprintf(`{"product":%d,"target_start":"%s","target_end":"%s","engagement_type":"%s","name":"%s"}`, id, target_date, target_date, e_type, name)
	expected := bytes.NewBuffer([]byte(json))

	assert.Nil(t, err)
	assert.NotNil(t, bytez)
	assert.Equal(t, expected, bytez)
}

func TestBuildMultipartFormBytez(t *testing.T) {
	k1 := "key"
	v1 := "value"
	k2 := "file"
	v2 := "path"
	fields := map[string]string{
		k1: v1,
		k2: v2,
	}
	data := []byte("test")

	c := defectdojo_client.NewDefectdojoClient("nil", "nil")
	bytez, header := c.BuildMultipartFormBytez(fields, data)

	// this isn't awesome, just asserting that those strings are set in the data
	// but not necessarily in the right order, should be close enough
	expected := `Content-Disposition: form-data; name="key"`
	expected2 := `Content-Disposition: form-data; name="file";`
	expected3 := `Content-Type: text/xml`

	assert.NotNil(t, bytez)
	assert.Contains(t, bytez.String(), expected)
	assert.Contains(t, bytez.String(), expected2)
	assert.Contains(t, bytez.String(), expected3)
	assert.Contains(t, header, "multipart/form-data; boundary=")
}

func TestDoPostCorrectlyBuildsUrl(t *testing.T) {
	api_path := "products"
	payload := defectdojo_client.Engagement{
		EngagementId: 5,
	}
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		full_url := fmt.Sprintf("/api/v2/%s/", api_path)
		assert.Equal(t, r.URL.Path, full_url)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	bytez, _ := c.BuildJsonRequestBytez(payload)
	resp, err := c.DoPost(api_path, bytez, defectdojo_client.APPLICATION_JSON)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDoGetCorrectlyBuildsUrl(t *testing.T) {
	api_path := "products"
	k1 := "param"
	v1 := "value"
	k2 := "param2"
	v2 := "value2"
	params := map[string]string{
		k1: v1,
		k2: v2,
	}
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		full_url := fmt.Sprintf("/api/v2/%s/", api_path)
		assert.Equal(t, r.URL.Path, full_url)
		assert.Equal(t, r.URL.Query().Get(k1), v1)
		assert.Equal(t, r.URL.Query().Get(k2), v2)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	c.DoGet(api_path, params)
}

func TestDoRequestReturnsResponse(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token api_key" {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "{}")
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	r, _ := http.NewRequest(http.MethodGet, mock_server.URL, nil)
	resp, err := c.DoRequest(r)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDoRequestComplainsifNot200Or201(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	req, _ := http.NewRequest(http.MethodGet, mock_server.URL, nil)
	_, err := c.DoRequest(req)

	assert.Nil(t, err)

	mock_server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req, _ = http.NewRequest(http.MethodGet, mock_server.URL, nil)
	_, err = c.DoRequest(req)

	assert.Nil(t, err)

	mock_server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req, _ = http.NewRequest(http.MethodGet, mock_server.URL, nil)
	_, err = c.DoRequest(req)

	assert.NotNil(t, err)
}

func TestDoRequestServerError(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token api_key" {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusTeapot)
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	r, _ := http.NewRequest(http.MethodGet, mock_server.URL, nil)
	resp, err := c.DoRequest(r)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "received status code of `418`")
	assert.Nil(t, resp)
}
