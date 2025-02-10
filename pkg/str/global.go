package str

import "strconv"

// StringToBool ...
func StringToBool(data string) bool {
	res, err := strconv.ParseBool(data)
	if err != nil {
		res = false
	}

	return res
}

// StringToInt ...
func StringToInt(data string) int {
	res, err := strconv.Atoi(data)
	if err != nil {
		res = 0
	}

	return res
}

// DefaultData ...
func DefaultData(data, defaultData string) string {
	if data == "" {
		return defaultData
	}

	return data
}
