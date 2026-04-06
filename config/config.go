package config

import (
	"fmt"
	"os"
)

func Config(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fmt.Sprintf("environment variable %q is not set", key)
	}

	return value
}
