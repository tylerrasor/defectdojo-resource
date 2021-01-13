package concourse

import (
	"io/ioutil"
)

func (w *Worker) ReadFile(path string) ([]byte, error) {
	w.LogDebug("checking for file: %s", path)
	bytez, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytez, nil
}
