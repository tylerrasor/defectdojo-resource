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

	w := concourse.AttachToWorker(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	if err := in.Get(w); err != nil {
		logrus.SetOutput(os.Stderr)
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}
