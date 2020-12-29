package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource"
)

func TestSourceValidate(t *testing.T) {
	source := resource.Source{
		DefectDojoUrl: "",
	}
	err := source.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required parameter `defectdojo_url` not supplied.")
}

func TestSourceValidateChecksForHttpOrHttps(t *testing.T) {
	source := resource.Source{
		DefectDojoUrl: "url-without-http.com",
	}

	err := source.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Please provide http(s):// prefix")

	source = resource.Source{
		DefectDojoUrl: "http://url-that-should-work.com",
	}

	err = source.Validate()
	assert.Nil(t, err)

	source = resource.Source{
		DefectDojoUrl: "https://url-that-should-work.com",
	}

	err = source.Validate()
	assert.Nil(t, err)
}
