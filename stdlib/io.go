package stdlib

import (
	"bufio"
	"fmt"
	"os"
)

var io = []mapEntry{
	entry("println", println,
		"Concatenates arguments and prints with a newline at the end",
	),
	entry("print", print,
		"Concatenates arguments and prints without a newline at the end",
	),
	entry("printf", printf,
		"Formats the first string using remaining arguments and prints",
	),
	entry("read", read,
		"Reads a line from the console. Throws error if fails",
		"Usage: (read)",
	),
}

func println(args ...interface{}) {
	fmt.Println(args...)
}

func printf(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}

func print(args ...interface{}) {
	fmt.Print(args...)
}

func read() string {
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return text[0 : len(text)-1] // ignore the '\n' char
}
