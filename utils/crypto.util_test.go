package utils_test

import (
	"crypto"
	_ "crypto/sha256"
	"fmt"
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

func TestHashing(t *testing.T) {
	foobar := "foobar"
	hash := NewEncoder(&foobar).Hash(crypto.SHA256)
	assert.Equal(t, fmt.Sprintf("%x", hash), "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2")
}

func TestHex(t *testing.T) {
	foobar := "foobar"
	hex := NewEncoder(&foobar).Hex(crypto.SHA256)
	assert.Equal(t, hex, "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2")

	hex = NewEncoder(&foobar).Hex(0)
	assert.Equal(t, hex, "666f6f626172")
}
