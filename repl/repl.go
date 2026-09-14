package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/enotinc/hmp/evaluator"
	"github.com/enotinc/hmp/lexer"
	"github.com/enotinc/hmp/parser"
)

const PROMTP string = " ~$ "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	for {
		fmt.Print(PROMTP)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Insect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", err))
	}
}
