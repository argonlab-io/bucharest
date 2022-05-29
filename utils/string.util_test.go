package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestArrayStringContain(t *testing.T) {
	arr := []string{"foobar"}
	foobarIndex, containFoobar := ArrayStringContians(arr, "foobar")
	assert.True(t, containFoobar)
	assert.Equal(t, foobarIndex, 0)
}

func TestArrayStringContainIndex(t *testing.T) {
	arr := []string{"foobar", "fozbaz"}
	fozbazIndex, containFozbaz := ArrayStringContians(arr, "fozbaz")
	assert.True(t, containFozbaz)
	assert.Equal(t, fozbazIndex, 1)
}

func TestArrayStringContainNotFound(t *testing.T) {
	arr := []string{"foobar"}
	fozbazIndex, containFozbaz := ArrayStringContians(arr, "fozbaz")
	assert.False(t, containFozbaz, false)
	assert.Equal(t, fozbazIndex, len(arr))
}
