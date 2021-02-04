package main

import (
	"os"

	"github.com/tylerrasor/defectdojo-resource/internal/check"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func main() {
	w := concourse.AttachToWorker(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	if err := check.Check(w); err != nil {
		w.LogError("%s", err)
		os.Exit(1)
	}
}
