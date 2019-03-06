package main

import (
	"fmt"
)

type Item struct {
	Id      string
	Type    string
	Classes []string
}

func (i *Item) hasClass(name string) bool {
	for _, cls := range i.Classes {
		if cls == name {
			return true
		}
	}
	return false
}

func (i *Item) String() string {
	return fmt.Sprintf("[Id]%s,[Type]%s,[Classed]%v", i.Id, i.Type, i.Classes)
}
