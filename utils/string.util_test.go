package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestArrayStringContain(t *testing.T) {
	arr := []string{"foobar"}
	foobarIndex, containFoobar := ArrayStringContians(arr, "foobar")
	assert.True(t,containFoobar)
	assert.Equal(t, foobarIndex, 0)
}
