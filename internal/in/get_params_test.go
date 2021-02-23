package in_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/in"
)

func TestGetParamsValidate(t *testing.T) {
	params := in.GetParams{}

	err := params.ValidateParams()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Required parameter `report_type` not supplied.")
}

func TestGetParamsValidateSucceeds(t *testing.T) {
	params := in.GetParams{
		ReportType: "something",
	}

	err := params.ValidateParams()
	assert.Nil(t, err)
}
