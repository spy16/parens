// Package repl provides facilities to build an interactive REPL using parens.
package repl

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spy16/parens"
	"github.com/spy16/parens/reader"
)

// New returns a new instance of REPL with given parens Env. Option values
// can be used to configure REPL input, output etc.
func New(env *parens.Env, opts ...Option) *REPL {
	repl := &REPL{
		rootEnv:   env,
		currentNS: func() string { return "" },
	}

	for _, option := range withDefaults(opts) {
		option(repl)
	}

	return repl
}

// REPL implements a read-eval-print loop for a generic Runtime.
type REPL struct {
	rootEnv     *parens.Env
	input       Input
	output      io.Writer
	mapInputErr ErrMapper
	currentNS   func() string
	factory     ReaderFactory

	banner      string
	prompt      string
	multiPrompt string

	printer Printer
}

// Input implementation is used by REPL to read user-input. See WithInput()
// REPL option to configure an Input.
type Input interface {
	SetPrompt(string)
	Readline() (string, error)
}

// Loop starts the read-eval-print loop. Loop runs until context is cancelled
// or input stream returns an irrecoverable error (See WithInput()).
func (repl *REPL) Loop(ctx context.Context) error {
	repl.printBanner()
	repl.setPrompt(false)

	for ctx.Err() == nil {
		err := repl.readEvalPrint()
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
	}

	return ctx.Err()
}

// readEval reads one form from the input, evaluates it and prints the result.
func (repl *REPL) readEvalPrint() error {
	forms, err := repl.read()
	if err != nil {
		switch err.(type) {
		case reader.Error:
			_ = repl.print(err)
		default:
			return err
		}
	}

	if len(forms) == 0 {
		return nil
	}

	res, err := parens.EvalAll(repl.rootEnv, forms)
	if err != nil {
		return repl.print(err)
	}
	if len(res) == 0 {
		return repl.print(nil)
	}

	return repl.print(res[len(res)-1])
}

func (repl *REPL) Write(b []byte) (int, error) {
	return repl.output.Write(b)
}

func (repl *REPL) print(v interface{}) error {
	if e, ok := v.(reader.Error); ok {
		s := fmt.Sprintf("%v (begin=%s, end=%s)", e, e.Begin, e.End)
		return repl.printer.Fprintln(repl.output, s)
	}
	return repl.printer.Fprintln(repl.output, v)
}

func (repl *REPL) read() ([]parens.Any, error) {
	var src string
	lineNo := 1

	for {
		repl.setPrompt(lineNo > 1)

		line, err := repl.input.Readline()
		err = repl.mapInputErr(err)
		if err != nil {
			return nil, err
		}

		src += line + "\n"

		if strings.TrimSpace(src) == "" {
			return nil, nil
		}

		rd := repl.factory.NewReader(strings.NewReader(src))
		rd.File = "REPL"

		form, err := rd.All()
		if err != nil {
			if errors.Is(err, reader.ErrEOF) {
				lineNo++
				continue
			}

			return nil, err
		}

		return form, nil
	}
}

func (repl *REPL) setPrompt(multiline bool) {
	if repl.prompt == "" {
		return
	}

	nsPrefix := repl.currentNS()
	prompt := repl.prompt

	if multiline {
		nsPrefix = strings.Repeat(" ", len(nsPrefix)+1)
		prompt = repl.multiPrompt
	}

	repl.input.SetPrompt(fmt.Sprintf("%s%s ", nsPrefix, prompt))
}

func (repl *REPL) printBanner() {
	if repl.banner != "" {
		fmt.Println(repl.banner)
	}
}
