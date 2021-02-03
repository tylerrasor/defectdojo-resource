package out_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/out"
)

func TestDecodeToPutRequestThrowsErrorWhenUnexpectedKey(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {},
		"params": {},
		"unexpectedkey": {}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	get, err := out.DecodeToPutRequest(w)

	assert.NotNil(t, err)
	assert.Nil(t, get)
}

func TestDecodeToPutRequestWorks(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "also exists",
			"product_name": "provided"
		},
		"params": {
			"report_type": "ZAP Scan",
			"path_to_report": "reports/report.txt"
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	get, err := out.DecodeToPutRequest(w)

	assert.Nil(t, err)
	assert.NotNil(t, get)
}
