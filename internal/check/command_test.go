package check_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/check"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func TestCheckErrorsWhenUnexpectedKeysInConcourseRequest(t *testing.T) {
	req := `{ "key": "unexpected" }`
	mock_stdin := bytes.NewBuffer([]byte(req))
	w := concourse.AttachToWorker(
		mock_stdin,
		nil,
		nil,
		nil,
	)

	err := check.Check(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload: ")
}

func TestCheckTurnsOnDebugWhenParamSet(t *testing.T) {
	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "api_key",
			"product_name": "product_name",
			"debug": true
		},
		"version": { "engagement_id": "5" }
	}`))
	var mock_stderr bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		&mock_stderr,
		nil,
		[]string{},
	)

	check.Check(w)
	w.LogDebug("what")

	assert.Contains(t, mock_stderr.String(), "debug logging on")
}

func TestCheckErrorsWhenGetEngagementFails(t *testing.T) {
	req := `
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "api_key",
			"product_name": "product_name"
		},
		"version": { "engagement_id": "5" }
	}`
	mock_stdin := bytes.NewBuffer([]byte(req))
	w := concourse.AttachToWorker(
		mock_stdin,
		nil,
		nil,
		nil,
	)

	err := check.Check(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting engagement: ")
}

func TestCheckPutsEngagementIdAsVersion(t *testing.T) {
	id := 15
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		e := fmt.Sprintf(`{ "id":%d,"target_start":"2021-02-04","target_end":"2021-02-04","product":5,"name":"name"}`, id)
		io.WriteString(w, e)
	}))

	req := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "api_key",
			"product_name": "product_name"
		},
		"version": {
			"engagement_id": "%s"
		}
	}`, mock_server.URL, fmt.Sprint(id))
	mock_stdin := bytes.NewBuffer([]byte(req))
	var mock_stdout bytes.Buffer
	w := concourse.AttachToWorker(
		mock_stdin,
		nil,
		&mock_stdout,
		nil,
	)

	err := check.Check(w)

	assert.Nil(t, err)
	assert.NotNil(t, mock_stdout)

	var r []concourse.Version
	decoder := json.NewDecoder(&mock_stdout)
	decoder.Decode(&r)
	assert.Equal(t, fmt.Sprint(id), r[0].EngagementId)
}

func TestCheckGivesBunkVersionOnInitialCheckForNow(t *testing.T) {
	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "should not be making any network calls until this is actually implemented")
	}))

	req := fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "api_key",
			"product_name": "product_name"
		},
		"version": {}
	}`, mock_server.URL)
	mock_stdin := bytes.NewBuffer([]byte(req))
	var mock_stdout bytes.Buffer
	w := concourse.AttachToWorker(
		mock_stdin,
		nil,
		&mock_stdout,
		nil,
	)

	err := check.Check(w)

	assert.Nil(t, err)
	assert.NotNil(t, mock_stdout)

	var r []concourse.Version
	decoder := json.NewDecoder(&mock_stdout)
	decoder.Decode(&r)
	assert.Equal(t, "https://github.com/tylerrasor/defectdojo-resource/issues/29", r[0].EngagementId)
}
