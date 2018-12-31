package decconv

import (
	"bytes"
	"fmt"
	"strings"
)

// Decode32 converts number from input into arbitrary precision decimal number with given precision and scale parameters
//   precision - amount of meaningful digits before and after dot
//   scale     - amount of meaningful digits after dot
func Decode32(precision, scale int, input []byte) (int32, error) {
	source := input

	if len(source) == 0 {
		return 0, fmt.Errorf("non-empty decimal number expected, got no data")
	}

	var negative bool
	if source[0] == '-' {
		negative = true
		source = source[1:]
	}

	// pass until . collecting data right into the result. Will pass leading zeroes as well
	var integral int32
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
			return 0, fmt.Errorf("decoding error: `%s` is not a decimal number", string(input))
		}
		passingZeroes = false
		integralCount++
		if integralCount > integralLimit {
			return 0, fmt.Errorf("overflow error on decoding %s, got %s digit while only %d are allowed in integral part",
				string(input),
				humanCount(integralCount),
				integralLimit,
			)
		}
		integral = integral*10 + int32(v-'0')
	}

	scaleCount := 0
	meaningful := 1
	multiplier := int32(10)
	var frac int32
	for _, v := range fraction {
		if v < '0' || v > '9' {
			return 0, fmt.Errorf("decoding error: `%s` is not a decimal number", string(input))
		}
		if v == '0' {
			multiplier *= 10
			meaningful++
			continue
		}
		frac = frac*multiplier + int32(v-'0')
		multiplier = 10
		scaleCount += meaningful
		if scaleCount > scale {
			return 0, fmt.Errorf("overflow error on decoding %s, got %s digit while only %d are allowed in fraction part",
				string(input),
				humanCount(scaleCount),
				scale,
			)
		}
		meaningful = 1
	}

	res := integral*pow32[scale] + frac*pow32[scale-scaleCount]
	if negative {
		return int32(-res), nil
	}
	return int32(res), nil
}

// Encode32 encodes given decimal number v with given scale into string
func Encode32(scale int, v int32) string {
	buf := &bytes.Buffer{}
	if int32(v) < 0 {
		buf.WriteByte('-')
		v = ^v + 1
	}
	integral := v / pow32[scale]
	_, _ = fmt.Fprintf(buf, "%d", integral)
	if scale == 0 {
		return buf.String()
	}
	fraction := v % pow32[scale]
	if fraction != 0 {
		format := fmt.Sprintf(".%%0%dd", scale)
		data := fmt.Sprintf(format, fraction)
		_, _ = buf.WriteString(strings.TrimRight(data, "0"))
	}
	return buf.String()
}
