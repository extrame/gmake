package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	rootOrNoSkip   = iota
	rootAndSkip    = iota
	notRootAndSkip = iota
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
func (g *Doc) Exec(waitingForWatch bool, nodependencies bool, selectors ...string) {
	selectStr := ""
	if len(selectors) > 0 {
		selectStr = selectors[0]
	}
	logrus.Infoln("try to execute by selector", selectors)
	selected := g.Select(selectStr)
	logrus.Debugf("selected (%d)", len(selected))
	if len(selectors) > 0 && len(selected) == 0 {
		fmt.Printf("selector %s doesn't existed\n", selectors[0])
		os.Exit(0)
	}
	ctx := &Context{wait: waitingForWatch, directivesInStack: make(map[int]bool)}
	ctx.load()
	var execChan = make(chan Result, len(selected))
	signalNo := 0
	for _, dir := range selected {
		logrus.Infoln("exec selector", dir)
		if nodependencies {
			dir.Exec(g, ctx, execChan, signalNo, rootAndSkip)
		} else {
			dir.Exec(g, ctx, execChan, signalNo, rootOrNoSkip)
		}
		signalNo++
	}
	for {
		select {
		case signal := <-execChan:
			logrus.Infof("No. %d is executed", signal.Serial)
			signalNo--
			if signalNo <= 0 && !waitingForWatch {
				return
			}
		}
	}
	// TODO
}

func (g *Doc) Select(selector string) Doc {
	if selector == "" {
		return Doc{(*g)[len(*g)-1]}
	}
	s := Selector(selector)
	return g.selectByItem(s)
}

func (g *Doc) selectByItem(selectItem *Item) Doc {
	d := make(Doc, 0)
	logrus.Infoln("selected by selector", selectItem)
	for _, item := range *g {
		logrus.Debug("compare with item", item.Name)
		if selectItem.Id != "" && item.Name.Id != selectItem.Id {
			continue
		}
		logrus.Debug("|same id|")
		if selectItem.Type != "" && item.Name.Type != selectItem.Type {
			continue
		}
		logrus.Debug("|same type|")
		var testLength = len(selectItem.Classes)
		for _, cls := range selectItem.Classes {
			if item.Name.hasClass(cls) {
				testLength--
			}
		}
		if testLength > 0 {
			continue
		}
		logrus.Debug("|same class")
		d = append(d, item)
	}
	return d
}
