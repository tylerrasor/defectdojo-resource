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

func TestGetEngagementReturnsErrorOnServerError(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.GetEngagement("id")

	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestGetEngagementForReportTypeGetsMostRecent(t *testing.T) {
	id := 15
	target_date := time.Now()
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// defectdojo stores engagements with autoincrement ids
		// when returning a list, it sends them back as a list in ascending order
		// thus, the last entry in the list is the most recent
		e1 := fmt.Sprintf(`{"id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id-1, target_date, target_date, app_id-1, report_type)
		e2 := fmt.Sprintf(`{"id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
		e := fmt.Sprintf(`{"count": 2, "results": [%s, %s]}`, e1, e2)
		io.WriteString(w, e)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := &defectdojo_client.Product{
		Id: app_id,
	}

	e, err := c.GetEngagementForReportType(p, report_type)

	assert.Nil(t, err)
	assert.Equal(t, id, e.EngagementId)
	assert.Equal(t, app_id, e.ProductId)
}

func TestGetEngagementForReportTypeJustReturnsErrorOnServerError(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.GetEngagementForReportType(&defectdojo_client.Product{}, "bunk report type")

	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestGetEngagementForReportTypeReturnsErrorOnBadPayload(t *testing.T) {
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "invalid json")
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := &defectdojo_client.Product{
		Id: app_id,
	}

	e, err := c.GetEngagementForReportType(p, report_type)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error decoding response: ")
	assert.Nil(t, e)
}

func TestCreateEngagementSetsTheRequestParamsCorrectly(t *testing.T) {
	target_date := time.Now()
	// because we can't match _exactly_ the timestamp
	target_date_substr := fmt.Sprintf("%d-%02d-%02d", target_date.Year(), target_date.Month(), target_date.Day())
	target_date_with_time := fmt.Sprintf("%d-%02d-%02d %02d:%02d", target_date.Year(), target_date.Month(), target_date.Day(), target_date.Hour(), target_date.Minute())
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Header.Get("Content-Type"), "application/json")
		var e *defectdojo_client.Engagement
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&e)
		assert.Equal(t, app_id, e.ProductId)
		assert.Equal(t, target_date_substr, e.StartDate)
		assert.Equal(t, target_date_substr, e.EndDate)
		assert.Equal(t, "CI/CD", e.EngagementType)
		name := fmt.Sprintf("%s - %s", report_type, target_date_with_time)
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

func TestCreateEngagementDoesCloseWhenCloseEngagementSet(t *testing.T) {
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

func TestCreateEngagementReturnsErrorWhenTheCreatePostServerError(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   5,
	}

	e, err := c.CreateEngagement(&p, "bunk report type", true)

	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestCreateEngagementReturnsErrorWhenCloseEngagementPostServerError(t *testing.T) {
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
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	p := defectdojo_client.Product{
		Name: "app",
		Id:   app_id,
	}

	e, err := c.CreateEngagement(&p, report_type, true)

	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestUploadReport(t *testing.T) {
	id := 18
	target_date := "2021-01-01"
	app_id := 5
	report_type := "report"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := fmt.Sprintf(`{ "engagement":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, app_id, report_type)
		io.WriteString(w, e)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.UploadReport(id, report_type, nil)

	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, e.EngagementIdFromUpload, id)
	assert.Equal(t, e.EngagementName, report_type)
	assert.Equal(t, e.ProductId, app_id)
	assert.Equal(t, e.StartDate, target_date)
	assert.Equal(t, e.EndDate, target_date)
}

func TestUploadReportReturnsErrorOnServerError(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	c := defectdojo_client.NewDefectdojoClient(mock_server.URL, "api_key")

	e, err := c.UploadReport(1, "bunk report type", []byte{})

	assert.Error(t, err)
	assert.Nil(t, e)
}
