package decconv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPowers10(t *testing.T) {
	max32 := pow32[len(pow32)-1]
	assert.Equal(t, int32(10), pow32[1])
	assert.Equal(t, int32(1000000000), max32)

	max64 := pow64[len(pow64)-1]
	assert.Equal(t, int64(10), pow64[1])
	assert.Equal(t, int64(1000000000)*int64(1000000000), max64)
}
