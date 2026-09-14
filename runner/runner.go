package runner

import (
	"fmt"
	"io"

	"github.com/enotinc/hmp/evaluator"
	"github.com/enotinc/hmp/lexer"
	"github.com/enotinc/hmp/object"
	"github.com/enotinc/hmp/parser"
)

func Run(input string, env *object.Enviroment, out io.Writer) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, fmt.Sprintf(" %s\n", evaluated.Inspect()))
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(out, fmt.Sprintf("\t%s\n", err))
	}
}
