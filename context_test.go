package bucharest_test

import (
	"context"
	"fmt"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/stretchr/testify/assert"
)

func XTestAddValuesToContext(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Fatalf("NewContext returned nil")
	}
	select {
	case x := <-ctx.Done():
		t.Errorf("<-c.Done() == %v want nothing (it should block)", x)
	default:
	}
	if got, want := fmt.Sprint(ctx), "context.Background"; got != want {
		t.Errorf("NewContext().String() = %q want %q", got, want)
	}

	key1 := 1
	value1 := "one"
	key2 := "two"
	value2 := 2

	ctx = AddValuesToContext(ctx, MapAny{
		key1: value1,
		key2: value2,
	})

	assert.Equal(t, value1, ctx.Value(key1))
	assert.Equal(t, value2, ctx.Value(key2))
}
