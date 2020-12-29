package models_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	models "github.com/tylerrasor/defectdojo-resource/models"
)

func TestSourceValidate(t *testing.T) {
	source := models.Source{
		DefectDojoUrl: "",
	}
	err := source.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required parameter `defectdojo_url` not supplied.")
}

func TestSourceValidateChecksForHttpOrHttps(t *testing.T) {
	source := models.Source{
		DefectDojoUrl: "url-without-http.com",
	}

	err := source.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Please provide http(s):// prefix")

	source = models.Source{
		DefectDojoUrl: "http://url-that-should-work.com",
	}

	err = source.Validate()
	assert.Nil(t, err)

	source = models.Source{
		DefectDojoUrl: "https://url-that-should-work.com",
	}

	err = source.Validate()
	assert.Nil(t, err)
}

func TestPutParamsValidate(t *testing.T) {
	params := models.PutParams{
		ReportType: "ZAP Scan",
	}

	err := params.Validate()
	assert.Nil(t, err)
}

func TestPutParamsValidateNoReportType(t *testing.T) {
	params := models.PutParams{}

	err := params.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required parameter `report_type` not supplied.")
}

func TestPutParamsValidateInvalidType(t *testing.T) {
	report_type := "invalid"
	params := models.PutParams{
		ReportType: report_type,
	}

	err := params.Validate()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", report_type)
	assert.EqualError(t, err, expected)
}

func TestPutParamsValidateNotYetImplemented(t *testing.T) {
	report_type := "Burp Scan"
	params := models.PutParams{
		ReportType: report_type,
	}

	err := params.Validate()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", report_type)
	assert.EqualError(t, err, expected)
}
