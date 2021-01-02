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

	assert.Error(t, err)
	message := fmt.Sprintf("product `%s` not found", name)
	assert.Equal(t, err.Error(), message)
	assert.Nil(t, p)
}

func TestGetProductReturnsAList(t *testing.T) {
	results := fmt.Sprintf(`{ "results": [ {}, {} ] }`)
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, results)
	}))
	defer mock_server.Close()
	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	name := "app"
	p, err := c.GetProduct(name)

	assert.Error(t, err)
	message := fmt.Sprintf("not sure how you did it, but got 2 results for product name `%s`", name)
	assert.Equal(t, err.Error(), message)
	assert.Nil(t, p)
}
