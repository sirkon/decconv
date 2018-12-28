package decconv

import (
	"fmt"
)

func humanCount(count int) string {
	switch count {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return fmt.Sprintf("%dth", count)
	}
}
