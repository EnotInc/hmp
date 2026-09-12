package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/enotinc/hmp/lexer"
	"github.com/enotinc/hmp/token"
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

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
