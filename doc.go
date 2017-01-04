package main

type Doc []*Directive

//Exec execute the Doc file
func (g *Doc) Exec(selectors ...string) {
	selectStr := ".main"
	if len(selectors) > 0 {
		selectStr = selectors[0]
	}
	selected := g.Select(selectStr)
	ctx := &Context{}
	for _, dir := range selected {
		dir.Exec(g, ctx)
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
