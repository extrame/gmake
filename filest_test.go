package main

import "testing"
import "fmt"

func TestFileFind(t *testing.T) {
	fmt.Print(IsFileChanged("../bireader/*.go"))
}
