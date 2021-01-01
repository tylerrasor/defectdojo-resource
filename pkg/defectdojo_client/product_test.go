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

func TestCreateEngagementSetsReportName(t *testing.T) {
	id := 18
	target_date := "2021-01-01"
	app_id := 5
	report_type := "report"
	mocK_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
		io.WriteString(w, e)
	}))

	c := defectdojo_client.NewDefectdojoClient(mocK_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   app_id,
	}
	e, err := c.CreateEngagement(&p, report_type)

	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, e.EngagementId, id)
	assert.Equal(t, e.ProductId, app_id)
	assert.Equal(t, e.EngagementName, report_type)
}
