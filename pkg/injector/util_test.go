package injector

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtil_Empty(t *testing.T) {
	assert := require.New(t)

	val := empty[int]()
	assert.Empty(val)
	val2 := empty[*int]()
	assert.Nil(val2)
	assert.Empty(val2)
}
