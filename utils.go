package main

import (
	"os"
	"strings"
)

func setEnv(name string, vals ...string) {
	os.Setenv(name, strings.Join(vals, getEnvSeperator()))
}

func appendEnv(name string, vals ...string) {
	oldEnv := os.Getenv(name)
	os.Setenv(name, oldEnv+getEnvSeperator()+strings.Join(vals, getEnvSeperator()))
}
