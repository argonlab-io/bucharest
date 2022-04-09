package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewEncoder(t *testing.T) {
	encoder := NewEncoder(nil)
	assert.NotNil(t, encoder)
	assert.Equal(t, encoder.Plaintext(), "")

	foobar := "foobar"
	encoder = NewEncoder(&foobar)
	assert.NotNil(t, encoder)
	assert.Equal(t, encoder.Plaintext(), foobar)
}

func TestRandomByteSucess(t *testing.T) {
	b := NewEncoder(nil).Random(1).Bytes()
	assert.NotEmpty(t, b)

	b = NewEncoder(nil).Random(0).Bytes()
	assert.Empty(t, b)
}
