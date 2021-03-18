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

func TestGetThrowsErrorWhenGetEngagementFails(t *testing.T) {
	name := "product_name"
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	var mock_stdin bytes.Buffer
	json := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "%s"
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
	id := 15
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == fmt.Sprintf("/api/v2/engagements/%d/", id) {
			w.WriteHeader(http.StatusOK)
			mock_date := "2021-03-17"
			e := fmt.Sprintf(`{ "id":%d,"target_start":"%s","target_end":"%s","product":5,"name":"name"}`, id, mock_date, mock_date)
			io.WriteString(w, e)
		} else {
			assert.Failf(t, "shouldn't call any other endpoints", "but called: %s", r.RequestURI)
		}
	}))

	var mock_stdin bytes.Buffer
	json := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "product_name"
		},
		"version": {
			"engagement_id": "%d"
		}
	}`, mock_server.URL, id)
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
