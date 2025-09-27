package injector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_GenerateServiceName(t *testing.T) {
	assert := require.New(t)

	type (
		test struct{} //nolint:unused
	)
	name := generateServiceName[test]()
	assert.Equal("injector.test", name)
	name = generateServiceName[*test]()
	assert.Equal("*injector.test", name)
}
