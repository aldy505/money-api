package decorator_test

import (
	"errors"
	"money-api/platform/decorator"
	"testing"
)

func TestErr(t *testing.T) {
	x := errors.New("hello world")
	wrapped := decorator.Err(x)

	if wrapped.Error() == x.Error() {
		t.Error("they must be different")
	}
}
