package config

import (
	"strconv"
)

func ConvertToInt(stringInput string) int {
	number, _ := strconv.ParseInt(stringInput, 10, 0)
	return int(number)
}

func ConvertToFloat(stringInput string) float64 {
	number, _ := strconv.ParseFloat(stringInput, 64)
	return float64(number)
}

func ConvertToBool(stringInput string) bool {
	done, _ := strconv.ParseBool(stringInput)
	return done
}
