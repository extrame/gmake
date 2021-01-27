package main

import (
	"os"
	"strings"
)

func setEnv(name string, vals ...string) string {
	newEnv := strings.Join(vals, getEnvSeperator())
	os.Setenv(name, newEnv)
	return newEnv
}

func appendEnv(name string, vals ...string) string {
	oldEnv := os.Getenv(name)
	newEnv := oldEnv + getEnvSeperator() + strings.Join(vals, getEnvSeperator())
	os.Setenv(name, newEnv)
	return newEnv
}
