package parens

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// NewREPL initializes a REPL session with given evaluator.
func NewREPL(exec Executor) *REPL {
	return &REPL{
		Exec:     exec,
		ReadIn:   defaultReadIn,
		WriteOut: defaultWriteOut,
	}
}

// Executor implementation is responsible for understanding
// and evaluating the input to generate a result.
type Executor interface {
	Execute(src string) (interface{}, error)
}

// REPL represents a session of read-eval-print-loop.
type REPL struct {
	Exec   Executor
	Banner string

	ReadIn   ReadInFunc
	WriteOut WriteOutFunc
}

// Start the REPL which reads from in and writes results to out.
func (repl *REPL) Start(ctx context.Context) error {
	if len(repl.Banner) > 0 {
		repl.WriteOut(repl.Banner, nil)
	}

	for {
		select {
		case <-ctx.Done():
			repl.WriteOut("Bye!", nil)
			return nil

		default:
			shouldExit := repl.readAndExecute()
			if shouldExit {
				repl.WriteOut("Bye!", nil)
				return nil
			}
		}
	}
}

func (repl *REPL) readAndExecute() bool {
	expr, err := repl.ReadIn()
	if err != nil {
		if err == io.EOF {
			return true
		}
		repl.WriteOut(nil, fmt.Errorf("read failed: %s", err))
		return false
	}

	if len(strings.TrimSpace(expr)) == 0 {
		return false
	}

	repl.WriteOut(repl.Exec.Execute(expr))
	return false
}

// ReadInFunc implementation is used by the REPL to read input.
type ReadInFunc func() (string, error)

// WriteOutFunc implementation is used by the REPL to write result.
type WriteOutFunc func(res interface{}, err error)

func defaultWriteOut(v interface{}, err error) {
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s\n", err)
	} else {
		fmt.Fprintln(os.Stdout, formatResult(v))
	}
}

func defaultReadIn() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	src, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// multiline source
	if strings.HasSuffix(src, "\\\n") {
		nl, err := defaultReadIn()
		if err != nil {
			return "", err
		}

		return strings.Trim(src, "\\\n") + "\n" + nl, nil
	}

	return strings.TrimSpace(src), nil
}
