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

func TestGetProductFound(t *testing.T) {
	id := 5
	name := "app"
	results := fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, name)
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, results)
	}))
	defer mock_server.Close()

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")
	p, err := c.GetProduct(name)

	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, p.Id, id)
}

func TestGetProductEmptyList(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{ "results": [ ] }`)
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	name := "app"
	p, err := c.GetProduct(name)

	message := fmt.Sprintf("product `%s` not found", name)
	assert.Errorf(t, err, message)
	assert.Nil(t, p)
}

func TestGetProductNotInList(t *testing.T) {
	id := 5
	name := "app"
	results := fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, name)
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, results)
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	other_name := "other name"
	p, err := c.GetProduct(other_name)

	message := fmt.Sprintf("product `%s` not found", other_name)
	assert.Errorf(t, err, message)
	assert.Nil(t, p)
}

func TestCreateEngagement(t *testing.T) {
	mocK_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	c := defectdojo_client.NewDefectdojoClient(mocK_server.URL, "api_key")

	p := defectdojo_client.Product{}
	_, err := c.CreateEngagement(&p)

	assert.Nil(t, err)
}
