GMake (1.0)
===========

一个用 Go 编写的轻量级构建工具。它读取项目中的 `GMakefile`，按其中以「类 CSS 选择器」
语法定义的指令执行命令，用于管理项目的构建、测试、清理等任务。

## 安装

    $ go install github.com/extrame/gmake@latest

在仓库内本地构建：

    $ go build -o gmake .

## 快速开始

在项目根目录创建 `GMakefile`（文件名区分大小写），定义一个指令：

    all {
        go build -o hello main.go
    }

执行：

    $ gmake all

`gmake` 会从当前目录开始向上查找 `GMakefile`，找到后切换到该文件所在目录执行。

### 命中的配置文件

每次执行都会把**实际命中的配置文件绝对路径**输出到 stderr，例如：

    $ gmake all
    gmake: using /home/me/project/GMakefile
    ...

这是判断「这次执行到底读了哪个文件」的唯一依据：向上查找时，离当前目录最近的
`GMakefile` 优先命中，子目录里的同名文件会遮蔽父目录的那个。

- **默认开启**，没有关闭开关，也不受 `--verbose`（只影响日志级别）或 `-nd`（只影响依赖）影响。
- 输出到 **stderr**，不会混进任务命令的 stdout，被其他脚本消费的输出不受污染。
- 在大小写不敏感的文件系统（如 macOS 默认的 APFS）上，会按磁盘上的**真实文件名**显示，
  因此 `Gmakefile` 和 `GMakefile` 能被区分出来。

### 执行规则

- 不带参数：执行文件中**最后一个**指令。
- 第一个参数是选择器：执行所有匹配的指令，多个匹配项按它们在 `GMakefile` 中出现的顺序执行。
- 选择器之后的参数是指令参数，可在命令中用 `$@`、`$1` 等占位符引用（见「7. 指令参数」）。

## 写法参考

### 1. 指令（target）

指令由「名称 + 可选的类 / ID / 属性 + 命令块」组成，与 CSS 选择器一致：

    build {                          # 纯名称
        go build ./...
    }

    build.deploy {                   # 类选择器
        echo deploying
    }

    build#main {                     # ID 选择器
        echo main
    }

    build[platform=linux] {          # 属性选择器
        echo linux build
    }

    build.production#main[platform=linux][arch=amd64] {   # 组合写法
        echo combined
    }

命名限制：

- 名称、类名、ID、属性名、属性值只允许字母、数字、`_`、`-`，且**必须以字母开头**；
  属性值以数字开头时值会丢失，导致属性选择器失效。
- 一个指令最多一个 ID，ID 重复会报错。

### 2. 命令

命令块内**每行一条命令**：

    all {
        go build -o hello main.go
        echo done
    }

注意事项：

- **不要写分号**。分号不是命令分隔符，而是会被拼进命令的最后一个参数，导致执行失败。
- **注释**：支持 `/* ... */` 块注释（可跨行，也可写在选择器或命令中间）；
  `#` 行注释**只在依赖括号 `(...)` 内有效**（`#` 后需为空白或行尾）。
  任务前、命令块内的 `#` 不是注释：指令前的 `#` 会直接报错，命令块内的 `#`
  会当作命令内容执行。`build#main` 里的 `#` 仍是 ID 选择器，
  双引号字符串内的内容不会被当作注释。
- gmake **不经过 shell**：命令按参数数组原样执行，管道 `|`、重定向 `>`、`&&`、
  通配符展开、`;` 都不会生效，这些字符会原样作为参数传给命令。因此即使参数里
  带有 `; rm -rf ...` 这样的片段，也只会成为某条命令的普通参数，不会被解释成
  新命令。这一点在让 AI 生成或传入参数时尤其重要。
- 唯一的语法特例是切换目录：`cd 目录 && 命令`，gmake 会把后续命令的工作目录切到该目录。
- 命令前加 `-` 表示该命令失败也不中断，继续执行后面的命令：

    all {
        -rm -rf dist
        mkdir dist
    }

- `$currentDir` 替换为 `GMakefile` 所在目录的绝对路径。
- `$变量名` 替换为 `var` 指令中定义的值；未定义时回退到同名环境变量。
- `$(命令)` 会先由 gmake 执行该命令，再用其输出替换当前位置（不经过 shell）。

### 3. 依赖

在指令名后用括号声明依赖，括号内写任意选择器：

    a.deploy { echo A }
    b.deploy { echo B }
    c (.deploy) { echo C }

`gmake c` 会先执行所有匹配 `.deploy` 的指令（a、b，按文件中的先后顺序），再执行 c。

多个依赖用逗号或换行分隔：

    lint {
        go vet ./...
    }

    c (.deploy, lint) {
        echo C
    }

- 依赖先于当前指令执行，当前指令等待所有依赖完成后再执行。
- 同一指令只应通过一条依赖路径引入：多个依赖选择器同时匹配到同一指令时，可能被重复执行。
- `gmake -nd c` 跳过所有依赖，只执行 c 本身。
- 依赖括号内可以用 `#` 写行注释：

    all (
        build.deploy, # 先构建
        # 再部署
        deploy.production
    ) {
        echo all
    }

### 4. 变量（var）

`var` 指令用 YAML 语法定义变量，通过类 / ID 供其他指令依赖：

    var.deploy {
        version: 1.0.0
        registry: git.example.com/demo
    }

    all.deploy (.deploy) {
        echo $version $registry
    }

- 变量的值里可以引用 `$currentDir` 和其他已定义变量。
- 变量会缓存到 `.gmake` 文件，用于判断变量是否发生变化。

#### 条件：变量更新时才执行

    all.deploy@version:updated (.deploy) {
        echo "version changed"
    }

`@变量名:updated` 表示仅当该变量相对上次运行（`.gmake` 缓存）发生变化时才执行该指令。

### 5. 环境变量（env）

`env` 指令用于设置进程环境变量，两种写法：

    env.deploy#GREET {          # 写法一：ID 指定变量名，命令内容为值
        hello
    }

    env.deploy {                # 写法二：YAML 批量设置
        GREET: hello
        MODE: dev
    }

在指令上追加 `.append` 类表示追加而不是覆盖：

    env.deploy#PATH.append {
        /opt/local/bin
    }

其他指令通过依赖 `.deploy` 引入这些环境变量：

    all.deploy (.deploy) {
        echo $GREET
    }

### 6. 执行选择器

选择器语法与 CSS 一致，支持类型、类、ID、属性的任意组合：

    $ gmake build
    $ gmake .deploy
    $ gmake #main
    $ gmake [platform=linux]
    $ gmake .production[platform=linux]

### 7. 指令参数

第一个命令行参数是选择器，其后的参数作为本次执行的指令参数，可在命令中用占位符引用：

    run {
        ./run.sh $@
    }

    $ gmake run a b

等价于执行 `./run.sh a b`。

- `$@`：展开为全部参数，每个参数仍保持为独立的命令参数；需独立成词，
  嵌在单词中间（如 `x$@`）时会合并成一个参数。
- `$1`、`$2`…：第 n 个参数，可与其他文本拼接（如 `--name=$1`）；越界时替换为空字符串。
- 未被占位符引用的参数会被忽略，不会自动附加到命令末尾。
- 参数属于本次执行，被执行的指令（含依赖）都可用占位符引用。
- 参数放在选择器之后；含空格的参数用引号包裹，如 `gmake run "hello world"`。
- 以 `-` 开头的参数也可正常传入（如 `gmake run --port 3000`），但必须写在选择器之后。

### 8. 完整示例

    var.deploy {
        bin: out.txt
    }

    env.deploy#GREET {
        hi
    }

    build.deploy (.deploy) {
        echo "building $bin"
        touch $bin
    }

    fmt {
        go fmt ./...
    }

    clean {
        -rm -f out.txt
    }

    all (build.deploy) {
        echo "build finished: $bin, $GREET"
    }

`gmake all` 的输出为：

    building out.txt
    build finished: out.txt, hi

## 目标清单（--list）

`gmake --list` 只解析 GMakefile，**不执行任何命令、不展开依赖**，退出码 0，
用来在执行前看清这个文件里有哪些目标：

    /abs/path/GMakefile

    build.deploy
      deps: .deploy
      commands: 2
        go build -o app
        echo done

`gmake --list --json` 输出同一份内容的结构化版本：

    {
      "file": "/abs/path/GMakefile",
      "targets": [
        {
          "selector": "build.deploy",
          "type": "build",
          "classes": ["deploy"],
          "id": "",
          "attributes": [],
          "conditions": [],
          "dependencies": [".deploy"],
          "commandCount": 2,
          "commands": [
            ["go", "build", "-o", "app"],
            ["echo", "done"]
          ]
        }
      ]
    }

- 字段是固定的，空值用 `[]` / `""` 表示，便于程序直接消费。
- `commands` 保留每条命令的参数切分，不拼成单个字符串，避免空格歧义。
- `var` / `env` 指令也会出现在清单里（`type` 分别为 `var`、`env`），它们是依赖链的一部分。
- 语法里**没有**「描述」这类信息，所以清单里也没有描述字段；想加目标说明只能靠命名，
  增加描述需要扩展语法，本次不做。
- `gmake: using <路径>` 仍然只写 stderr，stdout 是纯 JSON（需要更纯净时可用 `2>/dev/null`）。

## 命令行参数

    -nd           不执行依赖，只执行选中的指令
    -list, --list 列出所有目标，不执行任何命令
    -json, --json 与 --list 一起使用，输出 JSON
    -verbose N    日志级别：2 错误（默认）、3 警告、4 信息、5 调试
    -version      输出版本信息
    --help        Go flag 默认帮助

除以上开关外，每次执行都会固定向 stderr 输出 `gmake: using <配置文件绝对路径>`
（见「命中的配置文件」），该提示没有开关，`-verbose` / `-nd` 不影响它。

## 为什么适合交给 AI

- **不经过 shell**：命令以参数数组直接执行，AI 传入的参数不会变成 `; rm -rf` 这样的新命令；
  管道、重定向、`&&`、通配符都不生效，参数只是参数。
- **语法面小**：选择器 + 花括号 + 依赖括号，没有隐式规则；写错容易看出来，
  生成成功率高，出错时也容易定位。
- **可被程序发现**：`gmake --list --json` 让 AI 在执行前就知道有哪些目标、
  依赖什么、会跑哪些命令，不用读文件、不用猜。
- **输出可分离**：任务输出走 stdout，工具自身提示走 stderr，
  AI 消费 stdout 拿到的就是纯结果。

边界要如实说明：`$变量`、`$currentDir`、`$@`、`$1` 以及 `$(命令)` 都是 gmake 自己的替换，
其中 `$(命令)` 会被 gmake 直接执行（不经过 shell）；`cd 目录 && 命令` 是唯一的语法特例。
这些行为都发生在 GMakefile 文本里，命令行参数只做占位符替换，不会参与语法解析。

## 尚未实现

以下语法在代码中有入口但当前不生效，请勿依赖：

- `import` 指令
- `file` 指令
- `--watch` 文件变化监听

## License

见 [LICENSE](LICENSE)。
