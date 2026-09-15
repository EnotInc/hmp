package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/enotinc/hmp/object"
	"github.com/enotinc/hmp/runner"
)

const PROMTP string = "\033[36m ~$ \033[0m"
const header string = `
   █ █ █▀▀ █   █▀█   █▄▄▄█ █▀▀   █ ▄ █ █ ▀█▀ █ █ ▄
   █▀█ █▀▀ █   █▄█   █ █ █ █▀▀   █ █ █ █  █  █▀█ ▄
   ▀ ▀ ▀▀▀ ▀▀▀ █     ▀   ▀ ▀▀▀   ▀▀▀▀▀ ▀  ▀  ▀ ▀  *'help me please' repl
`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	fmt.Print(header)
	env := object.NewEnviroment()
	for {
		fmt.Print(PROMTP)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		runner.Run(line, env, out)
	}
}
