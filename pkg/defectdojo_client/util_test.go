package defectdojo_client_test

import (
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

	resp, err := c.DoPost(api_path, payload)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
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
