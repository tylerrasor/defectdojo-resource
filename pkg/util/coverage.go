package util

import (
	"fmt"
	"os"
	"testing"
)

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage/50123125#50123125
func FailIfCoverageLow(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.8 {
			fmt.Printf("Tests passed but coverage failed at %.2f (yes, that's a different number, it's calculated differently)\n", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
