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
			EngagementId: "0.1",
		},
	}

	err := w.OutputResponseToConcourse(r)

	assert.Nil(t, err)
	expected := fmt.Sprintf("{\"version\":{\"engagement_id\":\"%s\"}}\n", r.Version.EngagementId)
	assert.Equal(t, expected, mock_stdout.String())
}

func TestGetWorkDirReturnsCliArg(t *testing.T) {
	mock_args := []string{"arg0", "arg1", "arg2"}
	w := concourse.AttachToWorker(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		mock_args,
	)

	d := w.GetWorkDir()

	assert.Equal(t, d, "arg1")
}
