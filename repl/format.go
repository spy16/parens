package repl

import (
	"fmt"
	"reflect"
)

func formatResult(v interface{}) string {
	rval := reflect.ValueOf(v)
	return fmt.Sprintf("kind=%s, value=%v", rval.Kind(), rval)
}
