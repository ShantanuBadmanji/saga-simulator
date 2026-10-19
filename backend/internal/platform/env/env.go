package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrRequiredEnvNotSet = errors.New("required environment variable not set")

func lookupEnv(key envVar) (string, bool) {
	return os.LookupEnv(string(key))
}

func String(key envVar, fallback *string) (string, error) {
	value, ok := lookupEnv(key)
	if ok {
		return value, nil
	}

	if fallback != nil {
		return *fallback, nil
	}

	return "", fmt.Errorf("%q: %w", key, ErrRequiredEnvNotSet)
}

func Int(key envVar, fallback *int) (int, error) {
	value, ok := lookupEnv(key)
	if !ok {
		if fallback != nil {
			return *fallback, nil
		}
		return 0, fmt.Errorf("%q: %w", key, ErrRequiredEnvNotSet)
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%q: %w", key, err)
	}

	return n, nil
}

func Bool(key envVar, fallback *bool) (bool, error) {
	value, ok := lookupEnv(key)
	if !ok {
		if fallback != nil {
			return *fallback, nil
		}
		return false, fmt.Errorf("%q: %w", key, ErrRequiredEnvNotSet)
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%q: invalid boolean %q: %w", key, value, err)
	}

	return b, nil
}
