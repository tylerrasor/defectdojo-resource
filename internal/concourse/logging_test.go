package concourse

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogLevelSetToErrorByDefault(t *testing.T) {
	w := AttachToWorker(nil, nil, nil, nil)

	assert.Equal(t, w.logger.GetLevel(), logrus.ErrorLevel)
}

func TestEnableDebugLogSetsLevelCorrectly(t *testing.T) {
	var mock_stderr bytes.Buffer
	w := AttachToWorker(nil, &mock_stderr, nil, nil)

	w.LogDebug("test")

	assert.Equal(t, mock_stderr.String(), "")

	w.EnableDebugLog()

	assert.Equal(t, w.logger.GetLevel(), logrus.DebugLevel)
	assert.Contains(t, mock_stderr.String(), "debug logging on")
}
