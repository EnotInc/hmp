package main

import (
	"os"

	"github.com/enotinc/hmp/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
