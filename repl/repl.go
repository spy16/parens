package repl

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/c-bata/go-prompt"
)

// New initializes a REPL session with given evaluator.
func New(exec Executor, completer prompt.Completer) *REPL {
	if completer == nil {
		completer = func(prompt.Document) []prompt.Suggest {
			return nil
		}
	}

	return &REPL{
		exec:     exec,
		prompt:   ">> ",
		prompter: newPrompter(completer),
	}
}

// Executor implementation is responsible for understanding
// and evaluating the input to generate a result.
type Executor interface {
	Execute(src string) (interface{}, error)
}

// REPL represents a session of read-eval-print-loop.
type REPL struct {
	exec     Executor
	banner   string
	prompt   string
	prompter *prompt.Prompt
}

// Start the REPL which reads from in and writes results to out.
func (repl *REPL) Start(ctx context.Context, out io.Writer, errOut io.Writer) error {
	if len(repl.banner) > 0 {
		fmt.Fprintln(out, repl.banner)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			expr := repl.prompter.Input()
			if len(strings.TrimSpace(expr)) == 0 {
				continue
			}

			result, err := repl.exec.Execute(expr)
			if err != nil {
				fmt.Fprintf(errOut, "error: %s\n", err)
			} else {
				fmt.Fprintln(out, result)
			}
		}
	}
}

// SetBanner sets the message displayed at startup.
func (repl *REPL) SetBanner(banner string) {
	repl.banner = banner
}

func newPrompter(completer prompt.Completer) *prompt.Prompt {
	prm := prompt.New(nil, completer)
	return prm
}
