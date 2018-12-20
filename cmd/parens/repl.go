package main

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spy16/parens"
)

func newREPL(env parens.Scope) (*REPL, error) {
	ins, err := readline.New("> ")
	if err != nil {
		return nil, err
	}
	pr := &prompter{ins: ins}

	return &REPL{
		Env:      env,
		ReadIn:   pr.readIn,
		WriteOut: pr.writeOut,
	}, nil
}

// REPL represents a session of read-eval-print-loop.
type REPL struct {
	Env    parens.Scope
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

	repl.WriteOut(parens.Execute(strings.NewReader(expr), repl.Env))
	return false
}

// ReadInFunc implementation is used by the REPL to read input.
type ReadInFunc func() (string, error)

// WriteOutFunc implementation is used by the REPL to write result.
type WriteOutFunc func(res interface{}, err error)

type prompter struct {
	ins *readline.Instance
}

func (pr *prompter) readIn() (string, error) {
	src, err := pr.ins.Readline()
	if err != nil {
		if err == readline.ErrInterrupt {
			return "", io.EOF
		}
		return "", err
	}

	// multiline source
	if strings.HasSuffix(src, "\\\n") {
		nl, err := pr.readIn()
		return strings.Trim(src, "\\\n") + "\n" + nl, err
	}

	return strings.TrimSpace(src), nil
}

func (pr *prompter) writeOut(v interface{}, err error) {
	if err != nil {
		pr.ins.Write([]byte(fmt.Sprintf("error: %s\n", err)))
		return
	}
	pr.ins.Write([]byte(formatResult(v) + "\n"))
}

func formatResult(v interface{}) string {
	if v == nil {
		return "nil"
	}
	rval := reflect.ValueOf(v)
	switch rval.Kind() {
	case reflect.Func:
		return fmt.Sprintf("<function: %s>", rval.String())
	default:
		return fmt.Sprintf("%v", rval)
	}
}
