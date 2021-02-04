package out_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/out"
)

func TestPutParamsValidate(t *testing.T) {
	params := out.PutParams{
		ReportType: "ZAP Scan",
		ReportPath: "reports/report.txt",
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

func TestPutParamsValidateNoReportPath(t *testing.T) {
	params := out.PutParams{
		ReportType: "ZAP Scan",
	}

	err := params.ValidateParams()
	assert.EqualError(t, err, "Required parameter `path_to_report` not supplied.")
}
