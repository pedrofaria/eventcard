package utils

import "strconv"

func NumericToFloat32(v string) float32 {
	f, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0
	}

	return float32(f)
}
