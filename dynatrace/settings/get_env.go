package settings

import (
	"os"
	"strconv"
)

func GetIntEnv(name string, def, min, max int) int {
	sValue := os.Getenv(name)
	if len(sValue) == 0 {
		return def
	}
	iValue, err := strconv.Atoi(sValue)
	if err != nil {
		return def
	}
	if iValue < min || iValue > max {
		return def
	}
	return iValue
}
