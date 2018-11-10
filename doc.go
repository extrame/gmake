package main

import (
	"fmt"

	"github.com/golang/glog"
)

type Doc []*Directive

func (g Doc) Len() int {
	return len(g)
}

func (g Doc) Less(a int, b int) bool {
	first := g[a]
	second := g[b]
	if first.Name.Type == "env" {
		return true
	} else if first.Name.Type == "var" {
		if second.Name.Type == "env" {
			return false
		}
	}
	return false
}

func (g Doc) Swap(a, b int) {
	first := g[a]
	second := g[b]
	g[b] = first
	g[a] = second
}

//Exec execute the Doc file
func (g *Doc) Exec(waitingForWatch bool, selectors ...string) {
	selectStr := ".main"
	if len(selectors) > 0 {
		selectStr = selectors[0]
	}
	glog.Infoln("try to execute by selector", selectors)
	selected := g.Select(selectStr)
	glog.Infof("selected (%d)", len(selected))
	if len(selectors) > 0 && len(selected) == 0 {
		fmt.Printf("selector %s doesn't existed\n", selectors[0])
	}
	ctx := &Context{wait: waitingForWatch, directivesInStack: make(map[int]bool)}
	var execChan = make(chan Result, len(selected))
	signalNo := 0
	for _, dir := range selected {
		glog.Infoln("exec selector", dir)
		dir.Exec(g, ctx, execChan, signalNo)
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
	glog.Infoln("selected by selector", selectItem)
	for _, item := range *g {
		glog.Info("compare with item", item.Name)
		if selectItem.Id != "" && item.Name.Id != selectItem.Id {
			continue
		}
		glog.Info("|same id|")
		if selectItem.Type != "" && item.Name.Type != selectItem.Type {
			continue
		}
		glog.Info("|same type|")
		var testLength = len(selectItem.Classes)
		for _, cls := range selectItem.Classes {
			if item.Name.hasClass(cls) {
				testLength--
			}
		}
		if testLength > 0 {
			continue
		}
		glog.Info("|same class")
		d = append(d, item)
	}
	return d
}
