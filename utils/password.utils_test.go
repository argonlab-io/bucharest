package utils_test

import (
	"testing"

	. "github.com/argonlab-io/bucharest/utils"
	"github.com/stretchr/testify/assert"
)

func TestInvalidPassword(t *testing.T) {
	password := "password"
	hashedpassword, _ := HashPassword(password)
	assert.Equal(t, CheckPasswordHash("wrongpassword", hashedpassword), false)
}

func TestArgon2(t *testing.T) {
	password := "password"
	hashedpassword, _ := HashPasswordWithArgon2(password, nil)
	hashedpassword2, _ := HashPasswordWithArgon2(password, nil)
	checkCorrect, _ := CheckPasswordHashWithArgon2("password", hashedpassword)
	checkWorng, _ := CheckPasswordHashWithArgon2("wrongpassword", hashedpassword)
	assert.NotEqual(t, hashedpassword, hashedpassword2)
	assert.Equal(t, checkCorrect, true)
	assert.Equal(t, checkWorng, false)
}

func TestArgon2ErrorInvalidHash(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$foo$bar")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
	assert.ErrorIs(t, ErrInvalidHash, err)
}
