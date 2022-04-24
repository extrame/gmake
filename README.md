GMake (1.0)
===========

A very lightweight build utility written in Go. 

## Install

    $ go install github.com/extrame/gmake@latest

## Getting Started

### Basic Usage

In any project, create a file `gmakefile` and use the following syntax to create rules:

    target {
        command
    }

A more realistic example is shown below compiling a go program:

    all {
        go build -o hello main.go;
    }
    
    fmt {
        go fmt main.go;
    }
    
    clean {
        rm hello;
    }

### Define Target like css

You can define a target like css, for example, define a target with class

	clean.deploy {
	...
	}

Or define a target with ID

	clean#deployForTest{
	...
	}

### Dependency

You can make target depend on some target using CSS selector, like:

	deploy.main (.deploy) {
	...
	}

This target depend on all target with **deploy** class,if you have two target with **deploy** class before this target, they will execute before this target with text order.

For example

	a.deploy {
		echo 1
	}
	b.deploy {
		echo 2
	}
	c (.deploy){
		echo 3
	}

**gmake c** will print 1,2,3 in order

#### Default Action

if no target is specified, Gmake will execute the latest one.

### Variables

Use the following snippet for create an variable

```
var.deploy {
    test: git.oschina.net/mink-tech/bible
    test2: 23
}
```



#### Variable pseudo

You can use pseudo to follow the different status of variables

```
action@variable1:updated{
	echo 1
}
```

the echo will only executed when variabe1 is updated