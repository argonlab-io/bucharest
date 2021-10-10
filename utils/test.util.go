package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertPanic(t *testing.T, callback func(), expectPanic interface{}) {
	defer func() {
		r := recover()
		assert.Equal(t, expectPanic, r)
	}()
	callback()
}
