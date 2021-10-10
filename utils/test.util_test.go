package utils_test

import (
	"log"
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
)

func TestAssertPanic(t *testing.T) {
	AssertPanic(t, func() { log.Print("no panic") }, nil)
	AssertPanic(t, func() { log.Panic("ffffffff") }, "ffffffff")
}
