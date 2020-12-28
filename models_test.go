package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource"
)

func TestValidate(t *testing.T) {
	source := resource.Source{
		DefectDojoUrl: "",
	}
	err := source.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "Required parameter `DefectDojoUrl` not supplied.", err.Error())
}

func TestValidateChecksForHttp(t *testing.T) {
	source := resource.Source{
		DefectDojoUrl: "url-without-http.com",
	}

	err := source.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "Please provide http(s):// prefix", err.Error())
}
