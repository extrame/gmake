GMake (1.0)
===========

**English** | [简体中文](README.zh-cn.md)

A lightweight build tool written in Go. It reads a `GMakefile` from your project and runs
the commands defined there. Targets are declared with CSS-like selectors, and the same
selector syntax is used on the command line to pick what to run, which makes gmake handy
for builds, tests, and cleanup tasks.

## Install

    $ go install github.com/extrame/gmake@latest

To build from a local checkout:

    $ go build -o gmake .

## Getting started

Create a `GMakefile` in the project root (the filename is case-sensitive) and define a target:

    all {
        go build -o hello main.go
    }

Run it:

    $ gmake all

gmake walks up from the current directory until it finds a `GMakefile`, then changes to that
file's directory before executing the target.

### Which file was picked

On every run, gmake prints the **absolute path of the GMakefile it actually used** to stderr:

    $ gmake all
    gmake: using /home/me/project/GMakefile
    ...

This is the reliable way to answer "which file did this run read?". The lookup goes upward, so
the `GMakefile` closest to the current directory wins; a same-named file in a subdirectory
shadows the one above it.

- **Always on**, with no flag to disable it. It is unaffected by `--verbose` (log level only)
  and by `-nd` (dependency handling only).
- Written to **stderr**, so it never mixes into the stdout of your commands or pollutes output
  consumed by other scripts.
- On case-insensitive filesystems (such as the default APFS on macOS) it reports the **real
  filename on disk**, so `Gmakefile` and `GMakefile` can be told apart.

### Execution rules

- No arguments: runs the **last** target in the file.
- First argument is a selector: runs every matching target, in the order they appear in the file.
- Arguments after the selector are target arguments, referenced with `$@`, `$1`, and so on
  (see "7. Target arguments").

## Syntax reference

### 1. Targets

A target is a name, optional classes / ID / attributes, and a command block, matching CSS selectors:

    build {                          # plain name
        go build ./...
    }

    build.deploy {                   # class selector
        echo deploying
    }

    build#main {                     # ID selector
        echo main
    }

    build[platform=linux] {          # attribute selector
        echo linux build
    }

    build.production#main[platform=linux][arch=amd64] {   # combined
        echo combined
    }

Naming rules:

- Names, classes, IDs, attribute names, and attribute values may contain only letters, digits,
  `_`, and `-`, and **must start with a letter**. An attribute value that starts with a digit
  loses its value, which silently breaks the attribute selector.
- A target has at most one ID; duplicate IDs are an error.

### 2. Commands

Each line in a command block is **one command**:

    all {
        go build -o hello main.go
        echo done
    }

Things to keep in mind:

- **Do not use semicolons.** A semicolon is not a separator; it gets glued onto the last
  argument of the command and makes it fail.
- **Comments**: `/* ... */` block comments are supported anywhere (multi-line, and even inside
  a selector or a command). `#` line comments are supported **only inside dependency
  parentheses `(...)`** (the `#` must be followed by whitespace or end of line). A `#` before a
  target, or inside a command block, is not a comment: a leading `#` before a target is an
  error, and a `#` line inside a command block stays command content. The `#` in `build#main`
  is still an ID selector, and markers inside double quotes are never treated as comments.
- gmake **does not use a shell**: commands run with an argument vector, so pipes `|`,
  redirections `>`, `&&`, glob expansion, and `;` have no special meaning and are passed
  through as literal arguments. A value like `; rm -rf ...` stays an ordinary argument of one
  command instead of becoming a new command. This matters when an AI generates or passes
  arguments.
- The one syntactic exception is changing directory: `cd dir && command` makes gmake run the
  following commands with that working directory.
- Prefix a command with `-` to keep going when it fails:

    all {
        -rm -rf dist
        mkdir dist
    }

- `$currentDir` expands to the absolute path of the directory containing the `GMakefile`.
- `$name` expands to a variable declared by a `var` target; if it is not defined, gmake falls
  back to the environment variable of the same name.
- `$(command)` runs that command first (without a shell) and substitutes its output.

### 3. Dependencies

Declare dependencies in parentheses after the target name. Any selector is allowed inside:

    a.deploy { echo A }
    b.deploy { echo B }
    c (.deploy) { echo C }

`gmake c` runs every target matching `.deploy` first (a and b, in file order), then runs c.

Separate multiple dependencies with commas or newlines:

    lint {
        go vet ./...
    }

    c (.deploy, lint) {
        echo C
    }

- Dependencies run first; the current target waits for all of them to finish.
- A target should be reachable through a single dependency path: if several dependency
  selectors match the same target, it may run more than once.
- `gmake -nd c` skips all dependencies and runs only c.
- `#` line comments can be used inside the dependency parentheses:

    all (
        build.deploy, # build first
        # then deploy
        deploy.production
    ) {
        echo all
    }

### 4. Variables (var)

A `var` target declares variables in YAML syntax and exposes them through classes / IDs:

    var.deploy {
        version: 1.0.0
        registry: git.example.com/demo
    }

    all.deploy (.deploy) {
        echo $version $registry
    }

- Variable values may reference `$currentDir` and other variables that are already defined.
- Variables are cached in a `.gmake` file, which is used to detect changes.

#### Conditions: run only when a variable changed

    all.deploy@version:updated (.deploy) {
        echo "version changed"
    }

`@name:updated` means the target runs only when that variable differs from the previous run
(the `.gmake` cache).

### 5. Environment variables (env)

An `env` target sets environment variables for the process. There are two forms:

    env.deploy#GREET {          # form 1: ID names the variable, commands are its value
        hello
    }

    env.deploy {                # form 2: YAML, several variables at once
        GREET: hello
        MODE: dev
    }

Add the `.append` class to append instead of overwrite:

    env.deploy#PATH.append {
        /opt/local/bin
    }

Other targets pull these variables in by depending on `.deploy`:

    all.deploy (.deploy) {
        echo $GREET
    }

### 6. Running selectors

Selector syntax matches CSS and supports any combination of type, class, ID, and attributes:

    $ gmake build
    $ gmake .deploy
    $ gmake #main
    $ gmake [platform=linux]
    $ gmake .production[platform=linux]

### 7. Target arguments

The first command-line argument is the selector; the rest are target arguments that commands
can reference through placeholders:

    run {
        ./run.sh $@
    }

    $ gmake run a b

This is equivalent to running `./run.sh a b`.

- `$@` expands to all arguments, each still a separate command argument. It must stand alone as
  a word; embedded in a word (such as `x$@`) it merges into a single argument.
- `$1`, `$2`, ... are the n-th argument and can be concatenated with other text (for example
  `--name=$1`). Out-of-range indexes expand to an empty string.
- Arguments not referenced by a placeholder are ignored; they are not appended to commands.
- Arguments belong to the current run; every executed target (including dependencies) can use them.
- Put arguments after the selector. Quote arguments that contain spaces, e.g.
  `gmake run "hello world"`.
- Arguments starting with `-` also work (e.g. `gmake run --port 3000`), but must come after the
  selector.

### 8. Complete example

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

The output of `gmake all` is:

    building out.txt
    build finished: out.txt, hi

## Target listing (--list)

`gmake --list` only parses the GMakefile. It **runs no commands and does not expand
dependencies**, and exits with status 0. Use it to see what a file offers before running it:

    /abs/path/GMakefile

    build.deploy
      deps: .deploy
      commands: 2
        go build -o app
        echo done

`gmake --list --json` prints the same information in a structured form:

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

- Field names are stable, and empty values are `[]` / `""` so programs can consume the output
  directly.
- `commands` keeps each command split into arguments instead of joining it into one string,
  which avoids whitespace ambiguity.
- `var` and `env` targets are listed too (`type` is `var` or `env`); they are part of the
  dependency chain.
- The grammar has **no** notion of a description, so the listing has no description field.
  Documenting a target is limited to its name; adding descriptions would require new syntax,
  which is not part of this release.
- `gmake: using <path>` still goes to stderr, so stdout is pure JSON (add `2>/dev/null` if you
  want it even cleaner).

## Command-line options

    -nd           do not run dependencies, only the selected target
    -list, --list list all targets without running anything
    -json, --json print the target list as JSON (together with --list)
    -verbose N    log level: 2 error (default), 3 warning, 4 info, 5 debug
    -version      print version information
    --help        Go flag's default help

In addition to these flags, every run writes `gmake: using <absolute config path>` to stderr
(see "Which file was picked"). That message has no switch, and `-verbose` / `-nd` do not
affect it.

## Why it works well with AI

- **No shell**: commands execute as an argument vector, so arguments an AI passes cannot turn
  into new commands such as `; rm -rf`. Pipes, redirections, `&&`, and globs do nothing;
  arguments stay arguments.
- **Small grammar**: selectors, braces, and dependency parentheses, with no implicit rules.
  Mistakes are easy to spot, generation succeeds more often, and failures are easy to locate.
- **Discoverable by programs**: `gmake --list --json` tells an AI which targets exist, what they
  depend on, and which commands they run, before anything executes. No reading files, no guessing.
- **Separable output**: task output goes to stdout, the tool's own notices go to stderr, so an
  AI consuming stdout gets only the result.

The boundaries, stated plainly: `$name`, `$currentDir`, `$@`, `$1`, and `$(command)` are all
substitutions performed by gmake itself, and `$(command)` is executed directly by gmake (without
a shell). `cd dir && command` is the only syntactic exception. All of this happens inside the
GMakefile text; command-line arguments only fill placeholders and never take part in parsing.

## Not implemented

These constructs exist in the code but do not work yet; do not rely on them:

- the `import` target
- the `file` target
- `--watch` file watching

## License

See [LICENSE](LICENSE).
