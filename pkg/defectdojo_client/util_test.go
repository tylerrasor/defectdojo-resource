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

	assert.Errorf(t, err, "received status code of `500`")
	assert.Nil(t, resp)
}
