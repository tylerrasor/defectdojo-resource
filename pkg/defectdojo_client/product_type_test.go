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

func TestGetProductTypeReturnsNilWhenGetFails(t *testing.T) {
	product_name := "product name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	pt, err := c.GetProductType(product_name)

	assert.Error(t, err)
	assert.Nil(t, pt)
}

func TestGetProductTypeReturnsErrorWhenDecodeFails(t *testing.T) {
	product_name := "product name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pt := fmt.Sprintf(`{ bad json }`)
		io.WriteString(w, pt)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	product_type, err := c.GetProductType(product_name)

	assert.Error(t, err)
	assert.Nil(t, product_type)
	assert.Contains(t, err.Error(), "error decoding response: ")
}

func TestGetProductTypeReturnsErrorWhenProductNotFound(t *testing.T) {
	product_name := "product name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pt := fmt.Sprintf(`{ "valid json": "but did not find the right product_type" }`)
		io.WriteString(w, pt)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	product_type, err := c.GetProductType(product_name)

	assert.Error(t, err)
	assert.Nil(t, product_type)
	expected := fmt.Sprintf("product `%s` not found", product_name)
	assert.Equal(t, err.Error(), expected)
}

func TestGetProductTypeReturnsWeirdErrorWhenWeirdThingHappens(t *testing.T) {
	product_name := "product name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pt := fmt.Sprintf(`{ "results": [ { "name": "product 1" }, { "name": "product 2" } ] }`)
		io.WriteString(w, pt)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	product_type, err := c.GetProductType(product_name)

	assert.Error(t, err)
	assert.Nil(t, product_type)
	expected := fmt.Sprintf("not sure how you did it, but you got `%d` results for product_type name `%s`", 2, product_name)
	assert.Equal(t, err.Error(), expected)
}

func TestGetProductTypeReturnsCorrectProductType(t *testing.T) {
	product_name := "product name"
	product_id := 5
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		pt := fmt.Sprintf(`{ "results": [ { "name": "%s", "id": %d } ] }`, product_name, product_id)
		io.WriteString(w, pt)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	product_type, err := c.GetProductType(product_name)

	assert.Nil(t, err)
	assert.NotNil(t, product_type)
	assert.Equal(t, product_name, product_type.Name)
	assert.Equal(t, product_id, product_type.Id)
}
