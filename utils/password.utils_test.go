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

func TestArgon2ErrorNoVersion(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$argon2id$oversion$m=65536,t=3,p=2$Woo1mErn1s7AHf96ewQ8Uw$D4TzIwGO4XD2buk96qAP+Ed2baMo/KbTRMqXX00wtsU")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
}

func TestArgon2ErrorIncompatibleVersion(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$argon2id$v=88$m=65536,t=3,p=2$Woo1mErn1s7AHf96ewQ8Uw$D4TzIwGO4XD2buk96qAP+Ed2baMo/KbTRMqXX00wtsU")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
	assert.ErrorIs(t, ErrIncompatibleVersion, err)
}

func TestArgon2ErrorNoParam(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$argon2id$v=19$noparams$Woo1mErn1s7AHf96ewQ8Uw$D4TzIwGO4XD2buk96qAP+Ed2baMo/KbTRMqXX00wtsU")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
}

func TestArgon2ErrorSaltIllegalBase64(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$argon2id$v=19$m=65536,t=3,p=2$notbase64$D4TzIwGO4XD2buk96qAP+Ed2baMo/KbTRMqXX00wtsU")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
}

func TestArgon2ErrorHashIllegalBase64(t *testing.T) {
	valid, err := CheckPasswordHashWithArgon2("foobar", "$argon2id$v=19$m=65536,t=3,p=4$mGDVgi3ZS8iU6uJf_-8W3g$malformbase64")
	assert.Equal(t, valid, false)
	assert.Error(t, err)
}
