package defectdojo_client_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func TestGetEngagement(t *testing.T) {
	id := 18
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		mock_date := "2021-02-04"
		e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":5,"name":"name"}`, id, mock_date, mock_date)
		io.WriteString(w, e)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.GetEngagement(fmt.Sprint(id))

	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, id, e.EngagementId)
}

func TestCreateEngagementSetsTheRequestParamsCorrectly(t *testing.T) {
	target_date := time.Now()
	// because we can't match _exactly_ the timestamp
	target_date_substr := fmt.Sprintf("%d-%02d-%02d %02d:%02d", target_date.Year(), target_date.Month(), target_date.Day(), target_date.Hour(), target_date.Minute())
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		var e *defectdojo_client.Engagement
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&e)
		assert.Equal(t, app_id, e.ProductId)
		assert.Contains(t, e.StartDate, target_date_substr)
		assert.Contains(t, e.EndDate, target_date_substr)
		assert.Equal(t, "CI/CD", e.EngagementType)
		name := fmt.Sprintf("%s - %s", report_type, target_date_substr)
		assert.Contains(t, e.EngagementName, name)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   app_id,
	}
	c.CreateEngagement(&p, report_type, false)
}

func TestCreateEngagementDoesntCloseWhenCloseEngagementNotSet(t *testing.T) {
	id := 18
	target_date := "2021-02-04"
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// pass back an active engagement
		w.WriteHeader(http.StatusOK)
		e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
		io.WriteString(w, e)
		// fail if we call the `close engagement` endpoint
		path := fmt.Sprintf("/api/v2/engagements/%d/close/", id)
		assert.NotEqual(t, path, r.URL.Path)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   app_id,
	}
	c.CreateEngagement(&p, report_type, false)
}

func TestCreateEngagementDoesCloseWhenCloseEngagementNotSet(t *testing.T) {
	id := 18
	target_date := "2021-02-04"
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/engagements/" {
			w.WriteHeader(http.StatusOK)
			e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
			io.WriteString(w, e)
		} else {
			path := fmt.Sprintf("/api/v2/engagements/%d/close/", id)
			assert.Equal(t, path, r.URL.Path)
		}
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   app_id,
	}
	c.CreateEngagement(&p, report_type, true)
}

func TestUploadReport(t *testing.T) {
	id := 18
	target_date := "2021-01-01"
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
		io.WriteString(w, e)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.UploadReport(id, report_type, nil)

	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, e.EngagementId, id)
	assert.Equal(t, e.EngagementName, report_type)
	assert.Equal(t, e.ProductId, app_id)
	assert.Equal(t, e.StartDate, target_date)
	assert.Equal(t, e.EndDate, target_date)
}
