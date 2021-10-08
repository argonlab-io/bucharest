package bucharest_test

import (
	"fmt"
	"testing"

	. "github.com/argonlab-io/bucharest"
)

func TestNewContextWithOptions(t *testing.T) {
	ctx := NewContextWithOptions(nil)
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "bucharest.BuchatrestContext"; got != want {
		t.Errorf("NewContextWithOptions().String() = %q want %q", got, want)
	}
}
