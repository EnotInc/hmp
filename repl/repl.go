package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/enotinc/hmp/object"
	"github.com/enotinc/hmp/runner"
)

const PROMTP string = " ~$ "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

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
