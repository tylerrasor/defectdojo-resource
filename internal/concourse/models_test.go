package concourse_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/concourse"
)

func TestSourceValidate(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "",
		ApiKey:        "something",
		ProductName:   "app",
	}
	err := source.ValidateSource()
	assert.Error(t, err)
	assert.EqualError(t, err, "Required `defectdojo_url` not supplied.")
}

func TestSourceValidateChecksForHttpOrHttps(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "url-without-http.com",
		ApiKey:        "something",
		ProductName:   "app",
	}

	err := source.ValidateSource()
	assert.Error(t, err)
	assert.EqualError(t, err, "Please provide http(s):// prefix in `defectdojo_url`.")

	source = concourse.Source{
		DefectDojoUrl: "http://url-that-should-work.com",
		ApiKey:        "something",
		ProductName:   "app",
	}

	err = source.ValidateSource()
	assert.Nil(t, err)

	source = concourse.Source{
		DefectDojoUrl: "https://url-that-should-work.com",
		ApiKey:        "something",
		ProductName:   "app",
	}

	err = source.ValidateSource()
	assert.Nil(t, err)
}

func TestSourceValidateApiKeyMissing(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "http://something",
		ProductName:   "app",
	}

	err := source.ValidateSource()
	assert.Error(t, err)
	assert.EqualError(t, err, "Required `api_key` not supplied.")
}

func TestSourceValidateProductNameMissing(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "http://something",
		ApiKey:        "something",
	}

	err := source.ValidateSource()
	assert.Error(t, err)
	assert.EqualError(t, err, "Required `product_name` not supplied.")
}

func TestSourceValidateOptionalCreateProductNotSetDoesntError(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl: "http://something",
		ProductName:   "app",
		ApiKey:        "something",
	}

	err := source.ValidateSource()
	assert.Nil(t, err)
}

func TestSourceValidateOptionalCreateProductSetDoesErrorWhenProductTypeNotGiven(t *testing.T) {
	source := concourse.Source{
		DefectDojoUrl:            "http://something",
		ProductName:              "app",
		ApiKey:                   "something",
		CreateProductIfNotExists: true,
	}

	err := source.ValidateSource()
	assert.Error(t, err)
	assert.EqualError(t, err, "Optional `create_product_if_not_exist` set and required `product_type` not supplied.")
}
