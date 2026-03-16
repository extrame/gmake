package main

import (
	"os"
	"testing"
)

func TestAttributeSelector(t *testing.T) {
	// 测试属性选择器 [for=target]
	_, tokens := Lexer("GMAKE", `
		deploy[for=target1] {
			echo "Deploy for target1"
		}
		deploy[for=target2] {
			echo "Deploy for target2"
		}
		build[platform=linux] {
			echo "Build for linux"
		}
	`)

	ds := Parse("GMAKE", tokens)

	// 测试选择器 [for=target1]
	selected := ds.Select("[for=target1]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for [for=target1], expected 1 but got %d", len(selected))
	}

	// 测试选择器 [for=target2]
	selected = ds.Select("[for=target2]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for [for=target2], expected 1 but got %d", len(selected))
	}

	// 测试选择器 [platform=linux]
	selected = ds.Select("[platform=linux]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for [platform=linux], expected 1 but got %d", len(selected))
	}

	// 测试不存在的属性选择器
	selected = ds.Select("[for=nonexistent]")
	if len(selected) != 0 {
		t.Errorf("wrong selected directive number for [for=nonexistent], expected 0 but got %d", len(selected))
	}
}

func TestAttributeSelectorWithClass(t *testing.T) {
	// 测试属性选择器与类选择器组合
	_, tokens := Lexer("GMAKE", `
		deploy.production[for=target1] {
			echo "Production deploy for target1"
		}
		deploy.staging[for=target1] {
			echo "Staging deploy for target1"
		}
		deploy.production[for=target2] {
			echo "Production deploy for target2"
		}
	`)

	ds := Parse("GMAKE", tokens)

	// 测试组合选择器 .production[for=target1]
	selected := ds.Select(".production[for=target1]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for .production[for=target1], expected 1 but got %d", len(selected))
	}

	// 测试组合选择器 .staging[for=target1]
	selected = ds.Select(".staging[for=target1]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for .staging[for=target1], expected 1 but got %d", len(selected))
	}

	// 测试组合选择器 .production[for=target2]
	selected = ds.Select(".production[for=target2]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for .production[for=target2], expected 1 but got %d", len(selected))
	}
}

func TestAttributeSelectorWithID(t *testing.T) {
	// 测试属性选择器与ID选择器组合
	_, tokens := Lexer("GMAKE", `
		deploy#main[for=target1] {
			echo "Main deploy for target1"
		}
		deploy#secondary[for=target1] {
			echo "Secondary deploy for target1"
		}
	`)

	ds := Parse("GMAKE", tokens)

	// 测试组合选择器 #main[for=target1]
	selected := ds.Select("#main[for=target1]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for #main[for=target1], expected 1 but got %d", len(selected))
	}

	// 测试组合选择器 #secondary[for=target1]
	selected = ds.Select("#secondary[for=target1]")
	if len(selected) != 1 {
		t.Errorf("wrong selected directive number for #secondary[for=target1], expected 1 but got %d", len(selected))
	}
}

// 新增测试：读取示例文件并运行测试
func TestExampleFileAttributeSelector(t *testing.T) {
	// 读取示例文件
	exampleFile := "examples/GMakefile.attribute.example"
	content, err := os.ReadFile(exampleFile)
	if err != nil {
		t.Skipf("示例文件 %s 不存在，跳过测试: %v", exampleFile, err)
	}

	// 解析示例文件
	_, tokens := Lexer("GMAKE", string(content))
	ds := Parse("GMAKE", tokens)

	// 测试基本属性选择器
	t.Run("BasicAttributeSelectors", func(t *testing.T) {
		// 测试 [for=production]
		selected := ds.Select("[for=production]")
		if len(selected) != 2 { // 应该匹配 deploy[for=production] 和 deploy[for=production][region=us-east]
			t.Errorf("expected 2 directives for [for=production], got %d", len(selected))
		}

		// 测试 [for=staging]
		selected = ds.Select("[for=staging]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for [for=staging], got %d", len(selected))
		}

		// 测试 [platform=linux]
		selected = ds.Select("[platform=linux]")
		if len(selected) != 2 { // 应该匹配 build.production[platform=linux] 和 build.staging[platform=linux]
			t.Errorf("expected 2 directives for [platform=linux], got %d", len(selected))
		}

		// 测试 [platform=windows]
		selected = ds.Select("[platform=windows]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for [platform=windows], got %d", len(selected))
		}
	})

	// 测试组合选择器
	t.Run("CombinedSelectors", func(t *testing.T) {
		// 测试 .production[platform=linux]
		selected := ds.Select(".production[platform=linux]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for .production[platform=linux], got %d", len(selected))
		}

		// 测试 .staging[platform=linux]
		selected = ds.Select(".staging[platform=linux]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for .staging[platform=linux], got %d", len(selected))
		}

		// 测试 #unit[scope=fast]
		selected = ds.Select("#unit[scope=fast]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for #unit[scope=fast], got %d", len(selected))
		}

		// 测试 #integration[scope=slow]
		selected = ds.Select("#integration[scope=slow]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for #integration[scope=slow], got %d", len(selected))
		}
	})

	// 测试多个属性组合
	t.Run("MultipleAttributes", func(t *testing.T) {
		// 测试 [for=production][region=us-east]
		selected := ds.Select("[for=production][region=us-east]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for [for=production][region=us-east], got %d", len(selected))
		}

		// 测试 [for=production][region=eu-west]
		selected = ds.Select("[for=production][region=eu-west]")
		if len(selected) != 1 {
			t.Errorf("expected 1 directive for [for=production][region=eu-west], got %d", len(selected))
		}
	})

	// 测试不存在的选择器
	t.Run("NonExistentSelectors", func(t *testing.T) {
		// 测试不存在的属性
		selected := ds.Select("[for=nonexistent]")
		if len(selected) != 0 {
			t.Errorf("expected 0 directives for [for=nonexistent], got %d", len(selected))
		}

		// 测试不存在的组合
		selected = ds.Select(".nonexistent[platform=linux]")
		if len(selected) != 0 {
			t.Errorf("expected 0 directives for .nonexistent[platform=linux], got %d", len(selected))
		}
	})
}
