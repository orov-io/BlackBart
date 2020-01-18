package server

import (
	"os"
	"path/filepath"
	"runtime"
)

// DefaultPort is the returned port to listen if no other is provided.
const DefaultPort = "8080"

// Default service params
const (
	DefaultEnvValue = "UNKNOWN"
)

func envExist(envVar string) bool {
	return os.Getenv(envVar) != ""
}

// GetEnvPort searches for the port to start to serving stored on the portKey
// environment variable. It it not exists, return the DefaultPort constant.
// If this
func GetEnvPort(portKey string) string {
	if !envExist(portKey) {
		return DefaultPort
	}

	return os.Getenv(portKey)
}

// GetDir figures out what is the caller directory and returns dir as a path
// relative to caller
func GetDir(dir string) string {
	_, file, _, _ := runtime.Caller(1)
	FileDir := filepath.Dir(file)
	directory := filepath.Join(FileDir, dir)
	return directory
}

// GetEnvOrDefaultString tries to get the a value from the given key env variable.
// If nothing is retrieve, returns the DefaultEnvValue.
func GetEnvOrDefaultString(key string) string {
	if !envExist(key) {
		GetLogger().Warnf("Can't retrieve nothing from %v env variable", key)
		return DefaultEnvValue
	}

	return os.Getenv(key)
}

// EnvExist provides a quick way to know if a env variable is set.
func EnvExist(envVar string) bool {
	return envExist(envVar)
}
