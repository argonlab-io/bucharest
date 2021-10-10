package utils_test

import (
	"crypto"
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashSHA256Sucess(t *testing.T) {
	password := "foobar"
	expectedHash := "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2"
	actualHash := HashStringToHex(password, crypto.SHA256)
	assert.Equal(t, expectedHash, actualHash)
}

func TestRandomByteSucess(t *testing.T) {
	b := RandomBytes(make([]byte, 10))
	assert.NotEmpty(t, b)

}

func TestRandomByteNil(t *testing.T) {
	b := RandomBytes(nil)
	assert.Nil(t, b)
}
