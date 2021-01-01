package concourse_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func TestReadFileErrorWhenFileDoesNotExist(t *testing.T) {
	w := concourse.AttachToWorker(os.Stdin, os.Stderr, os.Stdout, os.Args)

	path := "path/to/your/file"
	err := w.FileExists(path)

	assert.Error(t, err)
	message := fmt.Sprintf("open %s: no such file or directory", path)
	assert.Equal(t, err.Error(), message)
}

func TestReadFileCorrectlyReturnsBytes(t *testing.T) {
	w := concourse.AttachToWorker(os.Stdin, os.Stderr, os.Stdout, os.Args)

	path := "/etc/hosts"
	err := w.FileExists(path)

	assert.Nil(t, err)
}
