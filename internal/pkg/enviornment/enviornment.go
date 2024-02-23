package environment

import "os"

func GetEnvVar(key, defaultValue string) string {
	val := os.Getenv(key)
	if len(val) < 1 {
		return defaultValue
	}
	return val
}
