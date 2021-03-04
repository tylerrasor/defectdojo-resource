package in_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/in"
)

func TestDecodeToGetRequestThrowsErrorWhenUnexpectedKey(t *testing.T) {
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
		nil,
	)

	get, err := in.DecodeToGetRequest(w)

	assert.NotNil(t, err)
	assert.Nil(t, get)
}

func TestDecodeToGetRequestWorks(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "must exist",
			"product_name": "also has to be here"
		},
		"params": {
			"report_type": "why did I make so many fields required"
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		nil,
	)

	get, err := in.DecodeToGetRequest(w)

	assert.Nil(t, err)
	assert.NotNil(t, get)
}

func TestDecodeToGetRequestThrowsErrorWhenSourceValidationFails(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something"
		},
		"params": {
			"report_type": "ZAP Scan"
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		nil,
	)

	get, err := in.DecodeToGetRequest(w)

	assert.Error(t, err)
	assert.Nil(t, get)
	assert.Contains(t, err.Error(), "invalid source config: ")
}

func TestDecodeToPutRequestThrowsErrorWhenParamsValidationFails(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "also exists",
			"product_name": "provided"
		},
		"params": {}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		nil,
		nil,
		nil,
	)

	get, err := in.DecodeToGetRequest(w)

	assert.Error(t, err)
	assert.Nil(t, get)
	assert.Contains(t, err.Error(), "invalid params config: ")
}
