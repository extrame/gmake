package main

import "testing"

func TestItem_hasClass(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		i    *Item
		args args
		want bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.hasClass(tt.args.name); got != tt.want {
				t.Errorf("Item.hasClass() = %v, want %v", got, tt.want)
			}
		})
	}
}
