package resource_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	resource "github.com/tylerrasor/defectdojo-resource/pkg"
)

func TestGetParamsValidate(t *testing.T) {
	params := resource.GetParams{}

	err := params.ValidateParams()
	assert.NotNil(t, err)
	assert.Errorf(t, err, "not implemented yet")
}
