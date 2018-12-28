/*
Here we store powers of 10 for uint32, uint64 and 128 bit unsigned integers
*/

package decconv

import (
	"github.com/sirkon/ds128"
	"math/big"
)

var pow32 []uint32
var pow64 []uint64
var pow128 []struct {
	lo uint64
	hi uint64
}

func init() {
	cur32 := uint32(1)
	for {
		pow32 = append(pow32, cur32)
		next32 := cur32 * 10
		if cur32 != next32/10 {
			break
		}
		cur32 = next32
	}

	cur64 := uint64(1)
	for {
		pow64 = append(pow64, cur64)
		next64 := cur64 * 10
		if cur64 != next64/10 {
			break
		}
		cur64 = next64
	}

	var lo uint64 = 1
	var hi uint64 = 0
	for {
		pow128 = append(pow128, struct {
			lo uint64
			hi uint64
		}{lo: lo, hi: hi})

		num := &big.Int{}
		num = num.SetUint64(hi)
		num.Lsh(num, 64)
		tmp := &big.Int{}
		tmp.SetUint64(lo)
		num.Add(num, tmp)
		num.Mul(num, big.NewInt(10))
		num.Rsh(num, 128)
		if num.Cmp(big.NewInt(0)) != 0 {
			break
		}

		nextLo, nextHi := ds128.Mul64(lo, hi, 10)

		lo = nextLo
		hi = nextHi
	}
}
