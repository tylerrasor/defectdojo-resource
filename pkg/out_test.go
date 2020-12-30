package resource_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource/pkg"
)

func TestDecodeFromOutSucceeds(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {},
		"params": {},
		"unexpectedkey": {}
	}`))

	out := resource.NewConcourse(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	req, err := resource.DecodeFromOut(out)

	assert.NotNil(t, err)
	assert.Nil(t, req)
}

func TestDecodeFromOutSetsErrorMessage(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {
			"defectdojo_url": "something"
		},
		"params": {
			"report_type": "something"
		}
	}`))

	out := resource.NewConcourse(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	req, err := resource.DecodeFromOut(out)

	assert.Nil(t, err)
	assert.NotNil(t, req)
}
