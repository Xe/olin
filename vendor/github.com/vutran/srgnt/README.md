# srgnt

> Simple command-line framework for Go

## Install

To get started, you need a working Go environment. Once available, grab the package here:

```bash
$ go get github.com/vutran/srgnt
```

## Usage

`srgnt` provides a very simple API. The following is an example of a minimal CLI app that has 1 command called "hello" which will print out "Hello, world!".

```go
package main

import (
    "flag"
    "fmt"
    "github.com/vutran/srgnt"
)

func Hello(_ *flag.FlagSet) {
    fmt.Println("Hello, world!")
}

func main() {
    cli := srgnt.CreateProgram("foo")
    cli.AddCommand("hello", Hello, "Prints \"Hello, world!\"")
    cli.Run()
}
```

By running `foo` in the example above, it will display the default help text that is generated for you automatically.

```bash
Usage:

        foo [flags] <command>

Commands:

        hello   Prints "Hello, world!"
```

### Flags and Commands

You can register as much flags and commands as you would like. Flags must precede commands.

The following example will output "Hello, world!":

```bash
$ foo hello
```

And with flags:

```bash
$ foo --name Vu hello
```

Our example above doesn't register any flags so the previous command can't read the `--name` flag. To register our flag, we need to update our `cli` instance from the original basic example.

```diff
package main

import (
    "flag"
    "fmt"
    "github.com/vutran/srgnt"
)

func Hello(flags *flag.FlagSet) {
+   name := flags.Lookup("name")
+   fmt.Printf("Hello, %s!\n", name.Value.String())
-   fmt.Println("Hello, world!")
}

func main() {
    cli := srgnt.CreateProgram("foo")
-   cli.AddCommand("hello", Hello, "Prints \"Hello, world!\"")
+   cli.AddCommand("hello", Hello, "Prints \"Hello, <name>!\"")
+   cli.AddStringFlag("name", "", "Set a name")
    cli.Run()
}
```

## License

MIT © [Vu Tran](https://github.com/vutran/srgnt)
