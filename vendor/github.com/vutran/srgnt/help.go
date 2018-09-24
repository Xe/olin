package srgnt

import (
	"bytes"
	"fmt"
	"github.com/vutran/ansi/colors"
	"github.com/vutran/ansi/styles"
)

func Help(p *Program) {
	var b bytes.Buffer

	for k, v := range p.Commands {
		fmt.Fprintf(&b, "\t%s\t%s\n", styles.Bold(colors.Cyan(k)), colors.White(v.Description))
	}

	fmt.Println(styles.Bold(colors.Yellow("Usage:")))
	fmt.Println("")
	fmt.Printf("\t%s [flags] <command>\n\n", styles.Bold(colors.Cyan(p.Name)))
	fmt.Println(styles.Bold(colors.Yellow("Commands:\n")))
	fmt.Println(b.String())
}
