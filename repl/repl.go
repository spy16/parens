package repl

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

// New initializes a REPL session with given evaluator.
func New(exec Executor) *REPL {
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
			expr, err := repl.ReadIn()
			if err != nil {
				repl.WriteOut(nil, fmt.Errorf("read failed: %s", err))
				continue
			}

			if len(strings.TrimSpace(expr)) == 0 {
				continue
			}

			repl.WriteOut(repl.Exec.Execute(expr))
		}
	}
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
	return reader.ReadString('\n')
}

func formatResult(v interface{}) string {
	rval := reflect.ValueOf(v)
	switch rval.Kind() {
	case reflect.Func:
		return fmt.Sprintf("func()")

	default:
		return fmt.Sprint(v)
	}
}
