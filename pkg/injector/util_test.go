package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_Empty(t *testing.T) {
	val := empty[int]()
	assert.Empty(t, val)
	val2 := empty[*int]()
	assert.Nil(t, val2)
	assert.Empty(t, val2)
}
