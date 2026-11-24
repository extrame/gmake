package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// listAttribute 是目标上的属性选择器，如 [platform=linux]
type listAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// listTarget 是一个指令的可序列化描述，只包含 GMakefile 语法里真实存在的信息
type listTarget struct {
	Selector     string          `json:"selector"`
	Type         string          `json:"type"`
	Classes      []string        `json:"classes"`
	ID           string          `json:"id"`
	Attributes   []listAttribute `json:"attributes"`
	Conditions   []string        `json:"conditions"`
	Dependencies []string        `json:"dependencies"`
	CommandCount int             `json:"commandCount"`
	Commands     [][]string      `json:"commands"`
}

type listOutput struct {
	File    string       `json:"file"`
	Targets []listTarget `json:"targets"`
}

// buildList 从已解析的文档生成目标清单，不执行任何命令
func buildList(file string, doc Doc) listOutput {
	out := listOutput{File: file, Targets: make([]listTarget, 0, len(doc))}
	for _, d := range doc {
		target := listTarget{
			Selector:     d.Name.Type + d.Name.String(),
			Type:         d.Name.Type,
			Classes:      append(make([]string, 0, len(d.Name.Classes)), d.Name.Classes...),
			ID:           d.Name.Id,
			Attributes:   make([]listAttribute, 0, len(d.Name.Attributes)),
			Conditions:   make([]string, 0, len(d.Name.Conditions)),
			Dependencies: append(make([]string, 0, len(d.Dependencies)), d.Dependencies...),
			CommandCount: len(d.Commands),
			Commands:     make([][]string, 0, len(d.Commands)),
		}
		for _, attr := range d.Name.Attributes {
			target.Attributes = append(target.Attributes, listAttribute{Name: attr.name, Value: attr.value})
		}
		for _, cond := range d.Name.Conditions {
			if s := cond.String(); s != "" {
				target.Conditions = append(target.Conditions, s)
			}
		}
		for _, c := range d.Commands {
			target.Commands = append(target.Commands, append(make([]string, 0, len(c.Parts)), c.Parts...))
		}
		out.Targets = append(out.Targets, target)
	}
	return out
}

// printListJSON 输出机器可读的目标清单
func printListJSON(w io.Writer, out listOutput) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

// printListText 输出人读的目标清单
func printListText(w io.Writer, out listOutput) {
	fmt.Fprintf(w, "%s\n", out.File)
	for _, t := range out.Targets {
		fmt.Fprintf(w, "\n%s\n", t.Selector)
		if len(t.Conditions) > 0 {
			fmt.Fprintf(w, "  conditions: %s\n", strings.Join(t.Conditions, " "))
		}
		if len(t.Dependencies) > 0 {
			fmt.Fprintf(w, "  deps: %s\n", strings.Join(t.Dependencies, ", "))
		}
		fmt.Fprintf(w, "  commands: %d\n", t.CommandCount)
		for _, parts := range t.Commands {
			fmt.Fprintf(w, "    %s\n", strings.Join(parts, " "))
		}
	}
}
