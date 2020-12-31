package resource_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource/internal"
)

// not sure how much value this actually provides, but at least we know we have
// a test around the expected output string to report back to concourse
func TestBuildRespone(t *testing.T) {
	var mock_stdout bytes.Buffer

	out := resource.NewConcourse(
		os.Stdin,
		os.Stderr,
		&mock_stdout,
		nil,
	)

	err := resource.OutputVersionToConcourse(out)

	assert.Nil(t, err)
	expected := "{\"version\":{\"version\":\"need to figure out unique combination of app name, version, build number, something\"}}\n"
	assert.Equal(t, expected, mock_stdout.String())
}
