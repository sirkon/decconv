package decconv

import (
	"bytes"
	"fmt"
	"github.com/sirkon/ds128"
	"math/big"
	"strings"
)

// Decode128 converts number from input into arbitrary precision decimal number with given precision and scale parameters
//   precision - amount of meaningful digits before and after dot
//   scale     - amount of meaningful digits after dot
func Decode128(precision, scale int, input []byte) (uint64, uint64, error) {
	source := input

	if len(source) == 0 {
		return 0, 0, fmt.Errorf("non-empty decimal number expected, got no data")
	}

	var negative bool
	if source[0] == '-' {
		negative = true
		source = source[1:]
	}

	// pass until . collecting data right into the result. Will pass leading zeroes as well
	var intLo uint64
	var intHi uint64
	var fraction []byte
	integralCount := 0
	passingZeroes := true
	integralLimit := precision - scale
	for i, v := range source {
		if passingZeroes {
			if v == '0' {
				continue
			}
		}
		if v == '.' {
			fraction = source[i+1:]
			break
		}
		if v < '0' || v > '9' {
			return 0, 0, fmt.Errorf("decoding error: `%s` is not a decimal number", string(input))
		}
		passingZeroes = false
		integralCount++
		if integralCount > integralLimit {
			return 0, 0, fmt.Errorf("overflow error on decoding %s, got %s digit while only %d are allowed in integral part",
				string(input),
				humanCount(integralCount),
				integralLimit,
			)
		}
		intLo, intHi = ds128.Mul64(intLo, intHi, 10)
		intLo, intHi = ds128.Add(intLo, intHi, uint64(v-'0'), 0)
	}

	scaleCount := 0
	meaningful := 1
	mulLo := uint64(10)
	var mulHi uint64
	var fracLo uint64
	var fracHi uint64
	for _, v := range fraction {
		if v < '0' || v > '9' {
			return 0, 0, fmt.Errorf("decoding error: `%s` is not a decimal number", string(input))
		}
		if v == '0' {
			mulLo = pow128[meaningful].lo
			mulHi = pow128[meaningful].hi
			meaningful++
			continue
		}
		fracLo, fracHi = ds128.Mul(fracLo, fracHi, mulLo, mulHi)
		fracLo, fracHi = ds128.Add(fracLo, fracHi, uint64(v-'0'), 0)
		mulLo = 10
		mulHi = 0
		scaleCount += meaningful
		if scaleCount > scale {
			return 0, 0, fmt.Errorf("overflow error on decoding %s, got %s digit while only %d are allowed in fraction part",
				string(input),
				humanCount(scaleCount),
				scale,
			)
		}
		meaningful = 1
	}

	resLo, resHi := ds128.Mul(intLo, intHi, pow128[scale].lo, pow128[scale].hi)
	fracLo, fracHi = ds128.Mul(fracLo, fracHi, pow128[scale-scaleCount].lo, pow128[scale-scaleCount].hi)
	resLo, resHi = ds128.Add(resLo, resHi, fracLo, fracHi)
	if negative {
		resLo, resHi = ds128.Negate(resLo, resHi)
	}
	return resLo, resHi, nil
}

func makeBig(lo uint64, hi uint64) *big.Int {
	v1 := &big.Int{}
	v1.SetUint64(hi)
	v1.Lsh(v1, 64)
	v2 := &big.Int{}
	v2.SetUint64(lo)
	return v1.Add(v1, v2)
}

// Encode128 convert a couple representing 128 bit decimal with given scale into string
func Encode128(scale int, lo, hi uint64) string {
	buf := &bytes.Buffer{}

	if int64(hi) < 0 {
		lo, hi = ds128.Negate(lo, hi)
		buf.WriteByte('-')
	}
	i := makeBig(lo, hi)
	f := &big.Int{}

	i, f = i.DivMod(i, makeBig(pow128[scale].lo, pow128[scale].hi), f)

	_, _ = fmt.Fprintf(buf, "%d", i)
	if scale == 0 {
		return buf.String()
	}
	if f.Cmp(big.NewInt(0)) != 0 {
		format := fmt.Sprintf(".%%0%dd", scale)
		data := fmt.Sprintf(format, f)
		_, _ = buf.WriteString(strings.TrimRight(data, "0"))
	}
	return buf.String()
}
