package in_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/internal/in"
)

func TestGetParamsValidate(t *testing.T) {
	params := in.GetParams{}

	err := params.ValidateParams()
	assert.NotNil(t, err)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "not implemented yet")
}
