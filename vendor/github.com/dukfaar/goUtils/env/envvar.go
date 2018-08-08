package env

import "os"

func GetDefaultEnvVar(key string, def string) string {
	result := os.Getenv(key)
	if result == "" {
		return def
	}

	return result
}
