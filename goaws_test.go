package goaws_test

import (
	"testing"

	"github.com/itispx/goaws"
)

func TestString(t *testing.T) {
	t.Parallel()

	str := "string-test"

	out := goaws.String(str)

	if out == nil {
		t.Error("expected non-nil pointer")
	} else if *out != str {
		t.Errorf("expected '%s', got '%s'", str, *out)
	}
}
