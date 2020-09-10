package parens_test

import (
	"testing"

	"github.com/spy16/parens"
)

func TestNew(t *testing.T) {
	p := parens.New()
	assertNotNil(t, p)
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Errorf("wanted non-nil value, got nil")
	}
}
