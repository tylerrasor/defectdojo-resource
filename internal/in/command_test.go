package in_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/in"
)

func TestGetThrowsErrorOnInvalidPayload(t *testing.T) {
	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(`
	{
		"source": {},
		"params": {},
		"unexpectedkey": {}
	}`))
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := in.Get(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload: ")
}

func TestGetTurnsOnDebugWhenParamSet(t *testing.T) {
	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "must exist",
			"product_name": "also must exist",
			"debug": true
		},
		"params": {
			"report_type": "needs to be here"
		},
		"version": {
			"engagement_id": "5"
		}
	}`))
	var mock_stderr bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		&mock_stderr,
		nil,
		[]string{},
	)

	in.Get(w)
	w.LogDebug("what")

	assert.Contains(t, mock_stderr.String(), "debug logging on")
}

func TestGetThrowsErrorWhenGetProductFails(t *testing.T) {
	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "must exist",
			"product_name": "also must exist"
		},
		"params": {
			"report_type": "needs to be here"
		},
		"version": {
			"engagement_id": "5"
		}
	}`))
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := in.Get(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting product: ")
}

func TestGetThrowsErrorWhenGetEngagementForReportTypeFails(t *testing.T) {
	name := "product_name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := 15
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", name) {
			w.WriteHeader(http.StatusOK)
			results := fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, name)
			io.WriteString(w, results)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))

	var mock_stdin bytes.Buffer
	json := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "%s"
		},
		"params": {
			"report_type": "report type"
		},
		"version": {
			"engagement_id": "5"
		}
	}`, mock_server.URL, name)
	mock_stdin.Write([]byte(json))

	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := in.Get(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting engagement: ")
}

func TestGetProperlySetsVersionOutput(t *testing.T) {
	name := "product_name"
	report_type := "report_type"
	id := 15
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", name) {
			w.WriteHeader(http.StatusOK)
			results := fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, name)
			io.WriteString(w, results)
		} else {
			w.WriteHeader(http.StatusOK)
			target_date := "2021-03-03"
			e1 := fmt.Sprintf(`{"id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id-1, target_date, target_date, id-1, report_type)
			e2 := fmt.Sprintf(`{"id":%d,"target_start":"%s","target_end":"%s","product":%d,"name":"%s"}`, id, target_date, target_date, id, "bunk report type")
			e := fmt.Sprintf(`{"count": 2, "results": [%s, %s]}`, e1, e2)
			io.WriteString(w, e)
		}
	}))

	var mock_stdin bytes.Buffer
	json := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "%s"
		},
		"params": {
			"report_type": "report type"
		},
		"version": {
			"engagement_id": "5"
		}
	}`, mock_server.URL, name)
	mock_stdin.Write([]byte(json))

	var mock_stdout bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		&mock_stdout,
		[]string{},
	)

	err := in.Get(w)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"version\":{\"engagement_id\":\"%d\"}}\n", id), mock_stdout.String())
}
