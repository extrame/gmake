package main

import (
	"github.com/golang/glog"
)

type Doc []*Directive

//Exec execute the Doc file
func (g *Doc) Exec(waitingForWatch bool, selectors ...string) {
	selectStr := ".main"
	if len(selectors) > 0 {
		selectStr = selectors[0]
	}
	selected := g.Select(selectStr)
	ctx := &Context{wait: waitingForWatch, directivesInStack: make(map[int]bool)}
	var execChan = make(chan Result, len(selected))
	signalNo := 0
	for _, dir := range selected {
		go dir.Exec(g, ctx, execChan, signalNo)
		signalNo++
	}
	for {
		select {
		case signal := <-execChan:
			glog.Infof("No. %d is executed", signal.Serial)
			signalNo--
			if signalNo <= 0 && !waitingForWatch {
				return
			}
		}
	}
	// TODO
}

func (g *Doc) Select(selector string) Doc {
	s := Selector(selector)
	return g.selectByItem(s)
}

func (g *Doc) selectByItem(selectItem Item) Doc {
	d := make(Doc, 0)
	for _, item := range *g {
		if selectItem.Id != "" && item.Name.Id != selectItem.Id {
			continue
		}
		if selectItem.Type != "" && item.Name.Type != selectItem.Type {
			continue
		}
		var testLength = len(selectItem.Classes)
		for _, cls := range selectItem.Classes {
			if item.Name.hasClass(cls) {
				testLength--
			}
		}
		if testLength > 0 {
			continue
		}
		d = append(d, item)
	}
	return d
}
