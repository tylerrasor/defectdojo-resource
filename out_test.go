package resource_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource"
)

func TestDecodeFromOutSucceeds(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {},
		"params": {},
		"unexpectedkey": {}
	}`))

	out := resource.NewOut(
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

	out := resource.NewOut(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	req, err := resource.DecodeFromOut(out)

	assert.Nil(t, err)
	assert.NotNil(t, req)
}

// not sure how much value this actually provides, but at least we know we have
// a test around the expected output string to report back to concourse
func TestBuildRespone(t *testing.T) {
	var mock_stdout bytes.Buffer

	out := resource.NewOut(
		os.Stdin,
		os.Stderr,
		&mock_stdout,
		nil,
	)

	err := resource.BuildResponse(out)

	assert.Nil(t, err)
	expected := "{\"version\":{\"version\":\"need to figure out unique combination of app name, version, build number, something\"}}\n"
	assert.Equal(t, expected, mock_stdout.String())
}
