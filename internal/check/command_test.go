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

func TestCheckErrorsWhenGetEngagementFails(t *testing.T) {
	req := `{ "source": {}, "version": {} }`
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

	req := fmt.Sprintf(`{ "source": { "defectdojo_url": "%s" }, "version": {} }`, mock_server.URL)
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

	var r concourse.Response
	decoder := json.NewDecoder(&mock_stdout)
	decoder.Decode(&r)
	assert.Equal(t, fmt.Sprint(id), r.Version.EngagementId)
}
