package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDirectiveParams 覆盖 $@ 与 $1..$n 的替换规则
func TestDirectiveParams(t *testing.T) {
	tests := []struct {
		name   string
		params []string
		origin []string
		want   []string
	}{
		{
			name:   "只有选择器没有参数",
			params: nil,
			origin: []string{"run", "$@"},
			want:   []string{"run"},
		},
		{
			name:   "独立$@展开为多个参数",
			params: []string{"a", "b c"},
			origin: []string{"run", "$@"},
			want:   []string{"run", "a", "b c"},
		},
		{
			name:   "位置参数可拼接且越界为空",
			params: []string{"a", "b c"},
			origin: []string{"--name=$1", "$2", "$9"},
			want:   []string{"--name=a", "b c", ""},
		},
		{
			name:   "位置参数从1开始",
			params: []string{"a"},
			origin: []string{"$0"},
			want:   []string{""},
		},
		{
			name:   "词内$@合并为单个参数",
			params: []string{"a", "b"},
			origin: []string{"x$@"},
			want:   []string{"xa b"},
		},
		{
			name:   "无参数时位置参数为空",
			params: nil,
			origin: []string{"run", "$1"},
			want:   []string{"run", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{Variables: make(map[string]interface{}), params: tt.params}
			_, parts := ctx.replaceVar(tt.origin...)
			if len(parts) != len(tt.want) {
				t.Fatalf("expected %v but got %v", tt.want, parts)
			}
			for i := range tt.want {
				if parts[i] != tt.want[i] {
					t.Errorf("arg %d expected %q but got %q", i, tt.want[i], parts[i])
				}
			}
		})
	}
}

// TestParamsWithVariables 覆盖参数与变量、$currentDir 共存
func TestParamsWithVariables(t *testing.T) {
	ctx := &Context{
		Variables: map[string]interface{}{"name": "demo"},
		params:    []string{"a"},
	}

	_, parts := ctx.replaceVar("$name", "$1", "$currentDir")
	want := []string{"demo", "a", getCurrentDirectory()}
	if len(parts) != len(want) {
		t.Fatalf("expected %v but got %v", want, parts)
	}
	for i := range want {
		if parts[i] != want[i] {
			t.Errorf("arg %d expected %q but got %q", i, want[i], parts[i])
		}
	}
}

// runInTempDir 在临时目录解析并执行 GMakefile，返回该目录路径
func runInTempDir(t *testing.T, content, selector string, params ...string) string {
	t.Helper()
	dir := t.TempDir()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	}()

	_, tokens := Lexer("GMAKE", content)
	ds := Parse("GMAKE", tokens)
	ds.Exec(false, false, append([]string{selector}, params...)...)
	return dir
}

func TestExecWithParams(t *testing.T) {
	t.Run("$@产生多个独立参数", func(t *testing.T) {
		dir := runInTempDir(t, `
			spread {
				touch $@
			}
		`, "spread", "one.txt", "two.txt")
		for _, name := range []string{"one.txt", "two.txt"} {
			if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
				t.Errorf("expected %s created: %v", name, err)
			}
		}
	})

	t.Run("$1只传入第一个参数", func(t *testing.T) {
		dir := runInTempDir(t, `
			single {
				touch file-$1.txt
			}
		`, "single", "x", "y")
		if _, err := os.Stat(filepath.Join(dir, "file-x.txt")); err != nil {
			t.Errorf("expected file-x.txt created: %v", err)
		}
		if _, err := os.Stat(filepath.Join(dir, "file-y.txt")); err == nil {
			t.Error("file-y.txt should not be created")
		}
	})

	t.Run("依赖指令也能取到参数", func(t *testing.T) {
		dir := runInTempDir(t, `
			inner {
				touch dep-$1.txt
			}
			outer (inner) {
				touch main-$1.txt
			}
		`, "outer", "v")
		for _, name := range []string{"dep-v.txt", "main-v.txt"} {
			if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
				t.Errorf("expected %s created: %v", name, err)
			}
		}
	})
}
