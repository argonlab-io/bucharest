package utils

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidHash = errors.New("the encoded hash is not in the correct format")
var ErrIncompatibleVersion = errors.New("incompatible version of argon2")

type Argon2HashingParam struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var DefaultArgon2HashingParams = &Argon2HashingParam{
	memory:      64 * 1024,
	iterations:  3,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPasswordWithArgon2(password string, params *Argon2HashingParam) (string, error) {
	if params == nil {
		params = DefaultArgon2HashingParams
	}
	p := params
	salt := NewEncoder(nil).Random(int(p.saltLength)).Bytes()

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	b64Salt := NewEncoder(nil).ReadBytes(salt).Base64()
	b64Hash := NewEncoder(nil).ReadBytes(hash).Base64()

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func CheckPasswordHashWithArgon2(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := getArgon2Params(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}

	return false, nil
}

func splitHashArgon2(encodedHash string) ([]string, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return vals, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, err
	}
	if version != argon2.Version {
		return nil, ErrIncompatibleVersion
	}

	return vals, nil
}

func getArgon2Params(encodedHash string) (p *Argon2HashingParam, salt, hash []byte, err error) {
	vals, err := splitHashArgon2(encodedHash)
	if err != nil {
		return nil, nil, nil, err
	}
	p = &Argon2HashingParam{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = NewDecoder(vals[4]).Bytes()
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = NewDecoder(vals[5]).Bytes()
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
