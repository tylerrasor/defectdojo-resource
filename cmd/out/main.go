package main

import (
	"os"

	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
	"github.com/tylerrasor/defectdojo-resource/internal/out"
)

func main() {
	w := concourse.AttachToWorker(
		os.Stdin,
		os.Stderr,
		os.Stdout,
		os.Args,
	)

	if err := out.Put(w); err != nil {
		w.LogError("%s", err)
		os.Exit(1)
	}
}
