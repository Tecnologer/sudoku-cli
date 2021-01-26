package utils

import "strconv"

//ParseStringArrayToIntArray converts string Array to int array
func ParseStringArrayToIntArray(input []string) (data []int, err error) {
	data = make([]int, len(input))

	for i, v := range input {
		if v == "" {
			v = "0"
		}

		data[i], err = strconv.Atoi(v)

		if err != nil {
			return
		}
	}

	return
}
