package concourse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func TestSourceValidate(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "",
		Username:      "something",
		ApiKey:        "something",
	}
	err := source.ValidateSource()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required `defectdojo_url` not supplied.")
}

func TestSourceValidateChecksForHttpOrHttps(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "url-without-http.com",
		Username:      "something",
		ApiKey:        "something",
	}

	err := source.ValidateSource()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Please provide http(s):// prefix in `defectdojo_url`.")

	source = concourse.Source{
		DefectDojoUrl: "http://url-that-should-work.com",
		Username:      "something",
		ApiKey:        "something",
	}

	err = source.ValidateSource()
	assert.Nil(t, err)

	source = concourse.Source{
		DefectDojoUrl: "https://url-that-should-work.com",
		Username:      "something",
		ApiKey:        "something",
	}

	err = source.ValidateSource()
	assert.Nil(t, err)
}

func TestSourceValidateUsernameMissing(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "http://something",
		ApiKey:        "something",
	}

	err := source.ValidateSource()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required `username` not supplied.")
}

func TestSourceValidateApiKeyMissing(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "http://something",
		Username:      "something",
	}

	err := source.ValidateSource()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required `api_key` not supplied.")
}
