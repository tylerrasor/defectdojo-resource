package resource_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource"
)

func TestPutParamsValidate(t *testing.T) {
	params := resource.PutParams{
		ReportType: "ZAP Scan",
	}

	err := params.Validate()
	assert.Nil(t, err)
}

func TestPutParamsValidateNoReportType(t *testing.T) {
	params := resource.PutParams{}

	err := params.Validate()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required parameter `report_type` not supplied.")
}

func TestPutParamsValidateInvalidType(t *testing.T) {
	report_type := "invalid"
	params := resource.PutParams{
		ReportType: report_type,
	}

	err := params.Validate()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", report_type)
	assert.EqualError(t, err, expected)
}

func TestPutParamsValidateNotYetImplemented(t *testing.T) {
	report_type := "Burp Scan"
	params := resource.PutParams{
		ReportType: report_type,
	}

	err := params.Validate()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", report_type)
	assert.EqualError(t, err, expected)
}
