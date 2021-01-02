package concourse

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func (w *Worker) ReadFile(path string) ([]byte, error) {
	logrus.New().Debugf("checking for file: %s", path)
	bytez, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytez, nil
}
