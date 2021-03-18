package out_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/out"
)

func TestPutThrowsErrorOnInvalidPayload(t *testing.T) {
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

	err := out.Put(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload: ")
}

func TestPutTurnsOnDebugWhenParamSet(t *testing.T) {
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
			"report_type": "needs to be here",
			"path_to_report": "somewhere?"
		}
	}`))
	var mock_stderr bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		&mock_stderr,
		nil,
		[]string{},
	)

	out.Put(w)
	w.LogDebug("what")

	assert.Contains(t, mock_stderr.String(), "debug logging on")
}

func TestPutChecksIfCreateProductFlagSetWhenProductLookupFails(t *testing.T) {
	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "must exist",
			"product_name": "also must exist"
		},
		"params": {
			"report_type": "needs to be here",
			"path_to_report": "somewhere?"
		}
	}`))
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := out.Put(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting product: ")
}

func TestPutTriesToCreateProductIfLookupFailsAndFlagSet(t *testing.T) {
	product_name := "product_name"
	product_type := "product_type"

	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "%s",
			"product_type": "%s",
			"create_product_if_not_exist": true
		},
		"params": {
			"report_type": "needs to be here",
			"path_to_report": "somewhere?"
		}
	}`, mock_server.URL, product_name, product_type)))
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := out.Put(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("error creating product '%s' for product_type '%s': ", product_name, product_type))
}

func TestPutDoesNotTryToCreateProductWhenProductLookupSuceeds(t *testing.T) {
	id := 15
	product_name := "product_name"

	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", product_name) {
			w.WriteHeader(http.StatusOK)
			results := fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, product_name)
			io.WriteString(w, results)
		} else if r.RequestURI != "/api/v2/engagements/" || r.Method != "POST" {
			assert.Fail(t, "shouldn't call any other endpoints", r.RequestURI)
		}
	}))

	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(fmt.Sprintf(`
	{
		"source": {
			"defectdojo_url": "%s",
			"api_key": "must exist",
			"product_name": "%s"
		},
		"params": {
			"report_type": "needs to be here",
			"path_to_report": "somewhere?"
		}
	}`, mock_server.URL, product_name)))
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		[]string{},
	)

	err := out.Put(w)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error getting or creating engagement: ")
}

func TestPutCorrectlyBuildsFullReportPathAndThrowsErrorIfFileReadFails(t *testing.T) {
	id := 15
	product_name := "product_name"

	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var results string
		w.WriteHeader(http.StatusOK)
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", product_name) {
			results = fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, product_name)
		} else if r.RequestURI == "/api/v2/engagements/" && r.Method == "POST" {
			results = fmt.Sprintf(`{ "id":%d,"target_start":"2021-03-07","target_end":"2021-03-07","product":5,"name":"name"}`, id)
		} else {
			assert.Fail(t, "shouldn't call any other endpoints", r.RequestURI)
		}
		io.WriteString(w, results)
	}))

	var mock_stdin bytes.Buffer
	file_path := "dir/file.ext"
	mock_stdin.Write([]byte(fmt.Sprintf(`
		{
			"source": {
				"defectdojo_url": "%s",
				"api_key": "must exist",
				"product_name": "%s",
				"debug": true
			},
			"params": {
				"report_type": "needs to be here",
				"path_to_report": "%s"
			}
		}`, mock_server.URL, product_name, file_path)))
	var mock_stderr bytes.Buffer
	workdir := "/tmp/random-string"
	w := concourse.AttachToWorker(
		&mock_stdin,
		&mock_stderr,
		nil,
		[]string{"name", workdir},
	)

	err := out.Put(w)

	assert.Contains(t, mock_stderr.String(), fmt.Sprintf("trying to read file: %s/%s", workdir, file_path))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error reading report file: ")
}

func TestPutThrowsErrorWhenReportUploadFails(t *testing.T) {
	id := 15
	product_name := "product_name"

	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var results string
		w.WriteHeader(http.StatusOK)
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", product_name) {
			results = fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, product_name)
		} else if r.RequestURI == "/api/v2/engagements/" && r.Method == "POST" {
			results = fmt.Sprintf(`{ "id":%d,"target_start":"2021-03-07","target_end":"2021-03-07","product":5,"name":"name"}`, id)
		} else if r.RequestURI == "/api/v2/import-scan/" && r.Method == "POST" {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			assert.Fail(t, "shouldn't call any other endpoints", r.RequestURI)
		}
		io.WriteString(w, results)
	}))

	// ok so this is hacky, but I need a "real" file to point at
	// so I'll just grab the temp file this test spawned and chop
	// up the absolute path to be recombined during file read
	full_path, _ := os.Executable()
	workdir := strings.Split(full_path, "out.test")[0]
	file_path := "out.test"

	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(fmt.Sprintf(`
		{
			"source": {
				"defectdojo_url": "%s",
				"api_key": "must exist",
				"product_name": "%s",
				"debug": true
			},
			"params": {
				"report_type": "needs to be here",
				"path_to_report": "%s"
			}
		}`, mock_server.URL, product_name, file_path)))
	var mock_stderr bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		&mock_stderr,
		nil,
		[]string{"name", workdir},
	)

	err := out.Put(w)

	assert.Contains(t, mock_stderr.String(), fmt.Sprintf("trying to read file: %s/%s", workdir, file_path))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error uploading report: ")
}

func TestPutWhenEverythingGoesRightOutputsVersionToConcourse(t *testing.T) {
	id := 15
	product_name := "product_name"

	mock_server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var results string
		w.WriteHeader(http.StatusOK)
		if r.RequestURI == fmt.Sprintf("/api/v2/products/?name=%s", product_name) {
			results = fmt.Sprintf(`{ "results": [ { "id": %d, "name": "%s" } ] }`, id, product_name)
		} else if r.RequestURI == "/api/v2/engagements/" && r.Method == "POST" {
			results = fmt.Sprintf(`{ "id":%d,"target_start":"2021-03-07","target_end":"2021-03-07","product":5,"name":"name"}`, id)
		} else if r.RequestURI == "/api/v2/import-scan/" && r.Method == "POST" {
			results = fmt.Sprintf(`{ "engagement":%d,"target_start":"2021-03-07","target_end":"2021-03-07","product":5,"name":"name"}`, id)
		} else {
			assert.Fail(t, "shouldn't call any other endpoints", r.RequestURI)
		}
		io.WriteString(w, results)
	}))

	// ok so this is hacky, but I need a "real" file to point at
	// so I'll just grab the temp file this test spawned and chop
	// up the absolute path to be recombined during file read
	full_path, _ := os.Executable()
	workdir := strings.Split(full_path, "out.test")[0]
	file_path := "out.test"

	var mock_stdin bytes.Buffer
	mock_stdin.Write([]byte(fmt.Sprintf(`
		{
			"source": {
				"defectdojo_url": "%s",
				"api_key": "must exist",
				"product_name": "%s"
			},
			"params": {
				"report_type": "needs to be here",
				"path_to_report": "%s"
			}
		}`, mock_server.URL, product_name, file_path)))
	var mock_stdout bytes.Buffer
	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		&mock_stdout,
		[]string{"name", workdir},
	)

	err := out.Put(w)

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"version\":{\"engagement_id\":\"%d\"}}\n", id), mock_stdout.String())
}
