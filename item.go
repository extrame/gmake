package main

type Item struct {
	Id         string
	Type       string
	Conditions []*condition
	// PSEUDOS    []string
	Classes []string
}

func (i *Item) String() string {
	var str string
	if i.Id != "" {
		str += "#" + i.Id
	}
	for _, v := range i.Classes {
		str += "." + v
	}
	for _, v := range i.Conditions {
		str += v.String()
	}
	return str
}

//now, only variable type is supported
type condition struct {
	name   string
	typ    string
	pseudo []string
}

func (c *condition) String() string {
	var str string
	for _, v := range c.pseudo {
		str += "$" + c.name + ":" + v
	}
	return str
}

func (c *condition) isNotAlready(ctx *Context) bool {
	if c.typ == "variable" {
		for _, p := range c.pseudo {
			switch p {
			case "updated":
				//when no change, not already, so return true
				return ctx.Variables[c.name] == ctx.oldVariables[c.name]
			}
		}
	}
	return false
}

func (i *Item) hasClass(name string) bool {
	for _, cls := range i.Classes {
		if cls == name {
			return true
		}
	}
	return false
}
