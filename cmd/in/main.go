package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	resource "github.com/tylerrasor/defectdojo-resource/pkg"
)

func main() {
	color.NoColor = false

	command := resource.NewConcourse(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	if err := command.Get(); err != nil {
		logrus.SetOutput(os.Stderr)
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
}
