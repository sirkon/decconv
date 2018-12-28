package decconv

import (
	"github.com/sirkon/ds128"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPowers10(t *testing.T) {
	max32 := pow32[len(pow32)-1]
	assert.Equal(t, uint32(10), pow32[1])
	assert.Equal(t, uint32(1000000000), max32)

	max64 := pow64[len(pow64)-1]
	assert.Equal(t, uint64(10), pow64[1])
	assert.Equal(t, uint64(1000000000)*uint64(1000000000)*10, max64)

	max128 := pow128[len(pow128)-1]
	assert.Equal(t, uint64(10), pow128[1].lo)
	assert.Equal(t, uint64(0), pow128[1].hi)
	max128lo, max128hi := ds128.Mul64(max64, 0, max64)
	assert.Equal(t, max128lo, max128.lo)
	assert.Equal(t, max128hi, max128.hi)
}
