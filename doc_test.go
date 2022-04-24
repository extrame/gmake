package main

import (
	"reflect"
	"testing"
)

func TestDoc_Exec(t *testing.T) {
	type args struct {
		selectors []string
	}
	tests := []struct {
		name string
		g    *Doc
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.Exec(false, false, tt.args.selectors...)
		})
	}
}

func TestDoc_Select(t *testing.T) {
	type args struct {
		selector string
	}
	tests := []struct {
		name string
		g    *Doc
		args args
		want Doc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Select(tt.args.selector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Doc.Select() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoc_selectByItem(t *testing.T) {
	type args struct {
		selectItem Item
	}
	tests := []struct {
		name string
		g    *Doc
		args args
		want Doc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.selectByItem(&tt.args.selectItem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Doc.selectByItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
