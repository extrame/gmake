package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func parseTestDoc(t *testing.T, content string) Doc {
	t.Helper()
	_, tokens := Lexer("GMAKE", content)
	doc := Parse("GMAKE", tokens)
	if tokens == nil {
		t.Fatal("failed to scan test GMakefile")
	}
	return doc
}

func TestBuildList(t *testing.T) {
	doc := parseTestDoc(t, `
		var.deploy {
			version: 1.0.0
		}
		build.production#main[platform=linux] (.deploy) {
			go build -o app
			echo done
		}
		test@version:updated {
			go test ./...
		}
	`)
	out := buildList("/tmp/GMakefile", doc)
	if out.File != "/tmp/GMakefile" {
		t.Errorf("wrong file, got %q", out.File)
	}
	if len(out.Targets) != 3 {
		t.Fatalf("wrong target number, expected 3 but got %d", len(out.Targets))
	}

	if got := out.Targets[0].Selector; got != "var.deploy" {
		t.Errorf("wrong selector, got %q", got)
	}
	if got := out.Targets[0].Classes; len(got) != 1 || got[0] != "deploy" {
		t.Errorf("wrong classes, got %v", got)
	}

	build := out.Targets[1]
	// selector 由类型名与 Item.String() 拼成，ID 在类之前，与内部表示一致
	if build.Selector != "build#main.production[platform=linux]" {
		t.Errorf("wrong selector, got %q", build.Selector)
	}
	if build.Type != "build" || build.ID != "main" {
		t.Errorf("wrong type or id, got %q %q", build.Type, build.ID)
	}
	if len(build.Attributes) != 1 || build.Attributes[0].Name != "platform" || build.Attributes[0].Value != "linux" {
		t.Errorf("wrong attributes, got %v", build.Attributes)
	}
	if len(build.Dependencies) != 1 || build.Dependencies[0] != ".deploy" {
		t.Errorf("wrong dependencies, got %v", build.Dependencies)
	}
	if build.CommandCount != 2 || len(build.Commands) != 2 {
		t.Errorf("wrong command count, got %d/%d", build.CommandCount, len(build.Commands))
	}
	if strings.Join(build.Commands[0], " ") != "go build -o app" {
		t.Errorf("wrong first command, got %v", build.Commands[0])
	}

	if len(out.Targets[2].Conditions) != 1 || out.Targets[2].Conditions[0] != "@version:updated" {
		t.Errorf("wrong conditions, got %v", out.Targets[2].Conditions)
	}
}

func TestBuildListEmptySlicesInJSON(t *testing.T) {
	doc := parseTestDoc(t, "simple {\n\techo hi\n}\n")
	var buf bytes.Buffer
	if err := printListJSON(&buf, buildList("/tmp/GMakefile", doc)); err != nil {
		t.Fatal(err)
	}
	got := buf.String()
	for _, field := range []string{`"classes": []`, `"attributes": []`, `"conditions": []`, `"dependencies": []`} {
		if !strings.Contains(got, field) {
			t.Errorf("expected empty array for %s in JSON, got:\n%s", field, got)
		}
	}

	// 输出必须能被程序直接反序列化
	var decoded listOutput
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("json output is not decodable: %v", err)
	}
	if len(decoded.Targets) != 1 || decoded.Targets[0].Selector != "simple" {
		t.Errorf("wrong decoded output, got %+v", decoded)
	}
}

func TestPrintListText(t *testing.T) {
	doc := parseTestDoc(t, `
		lint {
			go vet ./...
		}
		all (.deploy, lint) {
			echo all
		}
	`)
	var buf bytes.Buffer
	printListText(&buf, buildList("/tmp/GMakefile", doc))
	got := buf.String()
	for _, want := range []string{"/tmp/GMakefile", "lint", "all", "deps: .deploy, lint", "commands: 1", "echo all"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in text output, got:\n%s", want, got)
		}
	}
}
