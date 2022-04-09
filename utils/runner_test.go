package utils_test

import (
	"testing"
	"time"

	. "github.com/argonlab-io/bucharest/utils"
)

func TestRunUntil(t *testing.T) {
	RunUntil(func() bool { return true }, time.Second)
}
