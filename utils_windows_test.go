package main

import "testing"

func Test_getEnvSeperator(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnvSeperator(); got != tt.want {
				t.Errorf("getEnvSeperator() = %v, want %v", got, tt.want)
			}
		})
	}
}
