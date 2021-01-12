package pics

import (
	"os"
	"strconv"
)

func GetEnvInt(key string, def int) int {
	env := os.Getenv(key)
	value, err := strconv.Atoi(env)
	if err != nil {
		return def
	}
	return value
}

func GetEnvString(key string, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}

	return env
}
