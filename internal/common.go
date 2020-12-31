package resource

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func OutputVersionToConcourse(o *Concourse) error {
	version_str := "need to figure out unique combination of app name, version, build number, something"
	message := fmt.Sprintf("preparing to JSON encode response: %s", version_str)
	logrus.Debugln(message)

	version := Version{
		Version: version_str,
	}
	response := Response{
		Version: version,
	}

	return o.WriteResponseToConcourse(response)
}
