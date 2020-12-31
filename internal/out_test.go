package resource_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource/internal"
)

func TestDecodeToPutRequestThrowsErrorWhenUnexpectedKey(t *testing.T) {
	var mock_stdin bytes.Buffer

	mock_stdin.Write([]byte(`
	{
		"source": {},
		"params": {},
		"unexpectedkey": {}
	}`))

	c := resource.NewConcourse(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	get, err := resource.DecodeToPutRequest(c)

	assert.NotNil(t, err)
	assert.Nil(t, get)
}

func TestDecodeToPutRequestWorks(t *testing.T) {
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

	c := resource.NewConcourse(
		&mock_stdin,
		os.Stderr,
		os.Stdout,
		nil,
	)

	get, err := resource.DecodeToPutRequest(c)

	assert.Nil(t, err)
	assert.NotNil(t, get)
}
