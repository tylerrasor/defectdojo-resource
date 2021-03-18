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
		"version": {
			"engagement_id": "5"
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
		"version": {
			"engagement_id": "5"
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

func TestDecodeToGetRequestThrowsErrorWhenVersionNotProvided(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "http://something",
			"api_key": "also exists",
			"product_name": "provided"
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
	assert.Contains(t, err.Error(), "version did not have required `engagement_id`")
}
