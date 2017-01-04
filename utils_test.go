package main

import "testing"

func Test_setEnv(t *testing.T) {
	type args struct {
		name string
		vals []string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setEnv(tt.args.name, tt.args.vals...)
		})
	}
}

func Test_appendEnv(t *testing.T) {
	type args struct {
		name string
		vals []string
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appendEnv(tt.args.name, tt.args.vals...)
		})
	}
}
