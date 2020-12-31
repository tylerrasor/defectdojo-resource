package in_test

import (
	"bytes"
	"os"
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
		os.Stderr,
		os.Stdout,
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
			"defectdojo_url": "something"
		},
		"params": {}
	}`))

	w := concourse.AttachToWorker(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	get, err := in.DecodeToGetRequest(w)

	assert.Nil(t, err)
	assert.NotNil(t, get)
}
