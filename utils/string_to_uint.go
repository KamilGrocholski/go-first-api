package utils

import "strconv"

func StringToUint(value string) (uint, error) {
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		var r uint
		return r, err
	}

	out := uint(parsed)

	return out, nil
}
