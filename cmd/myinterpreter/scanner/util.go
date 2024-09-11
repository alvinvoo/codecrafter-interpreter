package scanner

import (
	"fmt"
	"math"
	"strconv"
)

func HandleNumberLiteral(v interface{}) interface{} {
	if num, ok := v.(float64); ok {
		// Separate the integer and fractional parts
		_, frac := math.Modf(num)

		// If the fractional part is not zero, it's a float
		if frac != 0 {
			return strconv.FormatFloat(num, 'f', -1, 64)
		} else {
			return fmt.Sprintf("%.1f", num)
		}
	}

	return v
}
