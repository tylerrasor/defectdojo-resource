package check_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/check"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func TestDecodeToCheckRequestThrowsErrorWhenUnexpectedKey(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {},
		"version": {},
		"unexpectedkey": {}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	check, err := check.DecodeToCheckRequest(w)

	assert.NotNil(t, err)
	assert.Nil(t, check)
}

func TestDecodeToCheckRequestWorks(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "api_key",
			"product_name": "product_name"
		},
		"version": {
			"engagement_id": "5"
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	check, err := check.DecodeToCheckRequest(w)

	assert.Nil(t, err)
	assert.NotNil(t, check)
}

func TestDecodeToCheckRequestWorksWhenNoVersionGiven(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "api_key",
			"product_name": "product_name"
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	check, err := check.DecodeToCheckRequest(w)

	assert.Nil(t, err)
	assert.NotNil(t, check)
}

func TestDecodeToCheckRequestThrowsErrorWhenSourceValidationFails(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": ""
		}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	check, err := check.DecodeToCheckRequest(w)

	assert.Error(t, err)
	assert.Nil(t, check)
	assert.Contains(t, err.Error(), "invalid source config: ")
}
