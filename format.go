package parens

import (
	"fmt"
	"reflect"
)

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
