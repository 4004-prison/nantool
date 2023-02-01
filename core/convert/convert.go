package convert

import (
	"errors"
	"strconv"
)

// StrArrToIntArr the string slice convert to the int slice
func StrArrToIntArr(strs []string) ([]int, error) {
	ints := make([]int, 0, len(strs))
	for _, s := range strs {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, errors.New("the slices contain non-integers")
		}
		ints = append(ints, int(i))
	}
	return ints, nil
}
