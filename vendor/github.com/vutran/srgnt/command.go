package srgnt

import (
	"flag"
)

type CommandFunction func(flags *flag.FlagSet)

type Command struct {
	Description string
	Callback    CommandFunction
}
