package main

import (
	"fmt"

	"github.com/extrame/goblet"
)

type Hooker struct {
}

func Get(ctx *goblet.Context) {
	fmt.Println("ok")
}
