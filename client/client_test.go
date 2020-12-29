package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/client"
)

func TestNewClientFailsWithBunkCreds(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer mock_server.Close()

	client, err := client.NewDefectdojoClient(mock_server.URL, "username", "", "bad_key")

	assert.NotNil(t, err)
	assert.Errorf(t, err, "error in authenticating to defectdojo instance at `%s`: ", mock_server.URL)
	assert.Nil(t, client)
}

// func TestTryToAuthSucceeds(t *testing.T) {
// 	source := models.Source{
// 		DefectDojoUrl: "http://something",
// 		Username:      "admin",
// 		ApiKey:        "good_key",
// 	}

// 	err := client.TryToAuth(&source)
// 	assert.Nil(t, err)
// }

// func TestTryToAuthFailsSetsCorrectErrorMessage(t *testing.T) {
// 	source := models.Source{
// 		DefectDojoUrl: "http://something",
// 		Username:      "admin",
// 		ApiKey:        "bad_key",
// 	}

// 	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "hello client")
// 	}))
// 	defer mock_server.Close()

// 	err := client.tryToAuth(&source)
// 	assert.NotNil(t, err)
// }
