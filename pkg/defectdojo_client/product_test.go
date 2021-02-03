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

func TestCreateProductReturnsErrorWhenPostFails(t *testing.T) {
	product_name := "product name"
	product_type := "product type"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusInternalServerError)
		}
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{ "results": [ { "id": 5 } ]}`)
		}
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p, err := c.CreateProduct(product_name, product_type)

	assert.Error(t, err)
	assert.Nil(t, p)
	assert.EqualError(t, err, "received status code of `500`")
}

func TestCreateProductReturnsErrorWhenResponseDecodeFails(t *testing.T) {
	product_name := "product name"
	product_type := "product type"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{ bad json }`)
		}
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, `{ "results": [ { "id": 5 } ]}`)
		}
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p, err := c.CreateProduct(product_name, product_type)

	assert.Error(t, err)
	assert.Nil(t, p)
	assert.Contains(t, err.Error(), "error decoding response: ")
}

func TestCreateProductReturnsCorrectProduct(t *testing.T) {
	product_name := "product name"
	product_type := "product type"
	product_type_id := 9
	product_id := 5
	product_desc := "required description"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			p_str := fmt.Sprintf(`{ "id": %d, "name": "%s", "prod_type": %d, "description": "%s" }`, product_id, product_name, product_type_id, product_desc)
			io.WriteString(w, p_str)
		}
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			pt_str := fmt.Sprintf(`{ "results": [ { "id": %d } ]}`, product_type_id)
			io.WriteString(w, pt_str)
		}
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p, err := c.CreateProduct(product_name, product_type)

	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, product_name, p.Name)
	assert.Equal(t, product_type_id, p.ProductTypeId)
	assert.Equal(t, product_id, p.Id)
	assert.Equal(t, product_desc, p.Description)
}
