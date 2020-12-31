package concourse_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

// not sure how much value this actually provides, but at least we know we have
// a test around the expected output string to report back to concourse
func TestOutputVersionResponseToConcourse(t *testing.T) {
	var mock_stdout bytes.Buffer

	w := concourse.AttachToWorker(
		os.Stdin,
		os.Stderr,
		&mock_stdout,
		nil,
	)

	r := concourse.Response{
		Version: concourse.Version{
			Version: "0.1",
		},
	}

	err := w.OutputResponseToConcourse(r)

	assert.Nil(t, err)
	expected := fmt.Sprintf("{\"version\":{\"version\":\"%s\"}}\n", r.Version.Version)
	assert.Equal(t, expected, mock_stdout.String())
}
