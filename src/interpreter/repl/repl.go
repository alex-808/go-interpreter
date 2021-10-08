package repl

import (
	"bufio"
	"fmt"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/lexer"
	"github.com/alex-davis-808/go-interpreter/src/interpreter/token"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	// Infinite while loop
	for {
		fmt.Printf(PROMPT)
		// read from line until encountering a newline
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Grab text from scanner and pass it into lexer
		line := scanner.Text()

		l := lexer.New(line)

		// iterate through tokens, printing one by one
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
