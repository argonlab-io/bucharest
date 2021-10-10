package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertPanic(t *testing.T, callback func(), expectPanic interface{}) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("The code did not panic")
		}
		assert.Equal(t, expectPanic, r)
	}()
	callback()
}
