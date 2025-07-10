package utils

import "fmt"

func BytesToOption(bytes []byte) (int, error) {
	if len(bytes) == 1 {
		b := bytes[0]
		// yeah baby, ascii codes
		if b > 48 && b < 57 {
			return int(b) - 48, nil
		}
		return -1, fmt.Errorf("error: cannot parse bytes to option")
	}
	return -1, fmt.Errorf("error: cannot parse bytes to option")
}
