package bucharest_test

import (
	"fmt"
	"testing"

	. "github.com/argonlab-io/bucharest"
	"github.com/argonlab-io/bucharest/utils"
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

func TestNewContextWithOptionsFromParentContext(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: NewContext()})
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

func TestNewContextWithOptionsWithNoAddtionalOptions(t *testing.T) {
	ctx := NewContextWithOptions(&ContextOptions{Parent: NewContext()})
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

	utils.AssertPanic(t, func() { ctx.ENV() }, ErrNoENV)
	utils.AssertPanic(t, func() { ctx.GORM() }, ErrNoGORM)
	utils.AssertPanic(t, func() { ctx.Log() }, ErrNoLogrus)
	utils.AssertPanic(t, func() { ctx.Redis() }, ErrNoRedis)
	utils.AssertPanic(t, func() { ctx.SQL() }, ErrNoSQL)
	utils.AssertPanic(t, func() { ctx.SQLX() }, ErrNoSQLX)
}
