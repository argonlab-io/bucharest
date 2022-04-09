package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
)

func TestAssertPanic(t *testing.T) {
	AssertPanic(t, func() {}, nil)
	AssertPanic(t, func() { panic("") }, "")
}
