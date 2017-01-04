package main

import "testing"

func Test_combiner(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combiner(tt.args.strs); got != tt.want {
				t.Errorf("combiner() = %v, want %v", got, tt.want)
			}
		})
	}
}
