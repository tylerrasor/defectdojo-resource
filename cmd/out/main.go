package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	resource "github.com/tylerrasor/defectdojo-resource"
)

func main() {
	color.NoColor = false

	command := resource.NewOut(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	err := command.Execute()
	if err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
