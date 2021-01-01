package defectdojo_client_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tylerrasor/defectdojo-resource/pkg/defectdojo_client"
)

func TestBuildAuthHeader(t *testing.T) {
	key := "api_key"

	k, v := defectdojo_client.BuildAuthHeader(key)

	assert.Equal(t, k, "Authorization")
	token_str := fmt.Sprintf("Token %s", key)
	assert.Equal(t, v, token_str)
}
