package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/in"
)

func main() {
	color.NoColor = false

	command := concourse.NewConcourse(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	if err := in.Get(command); err != nil {
		logrus.SetOutput(os.Stderr)
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}
