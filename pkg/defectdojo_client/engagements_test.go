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
