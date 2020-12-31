package out_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/out"
)

func TestPutParamsValidate(t *testing.T) {
	params := out.PutParams{
		ReportType: "ZAP Scan",
	}

	err := params.ValidateParams()
	assert.Nil(t, err)
}

func TestPutParamsValidateNoReportType(t *testing.T) {
	params := out.PutParams{}

	err := params.ValidateParams()
	assert.NotNil(t, err)
	assert.EqualError(t, err, "Required parameter `report_type` not supplied.")
}

func TestPutParamsValidateInvalidType(t *testing.T) {
	report_type := "invalid"
	params := out.PutParams{
		ReportType: report_type,
	}

	err := params.ValidateParams()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, is not a supported by Defectdojo (check that your format matches expected)", report_type)
	assert.EqualError(t, err, expected)
}

func TestPutParamsValidateNotYetImplemented(t *testing.T) {
	report_type := "Burp Scan"
	params := out.PutParams{
		ReportType: report_type,
	}

	err := params.ValidateParams()
	assert.NotNil(t, err)
	expected := fmt.Sprintf("The specified report type, `%s`, hasn't been implemented yet (pull requests welcome!)", report_type)
	assert.EqualError(t, err, expected)
}
