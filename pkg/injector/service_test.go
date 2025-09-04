package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateServiceName(t *testing.T) {
	type (
		test struct{} //nolint:unused
	)

	name := generateServiceName[test]()
	assert.Equal(t, "injector.test", name)
	name = generateServiceName[*test]()
	assert.Equal(t, "*injector.test", name)
}
