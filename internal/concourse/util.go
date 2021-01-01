package concourse

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func (w *Worker) FileExists(path string) error {
	logrus.New().Debugf("checking for file: %s", path)
	_, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return nil
}
