package main

import "testing"

func TestDirective_Exec(t *testing.T) {
	type args struct {
		doc *Doc
	}
	tests := []struct {
		name string
		d    *Directive
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.d.Exec(tt.args.doc, &Context{})
		})
	}
}
