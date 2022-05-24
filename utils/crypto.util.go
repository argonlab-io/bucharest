package utils

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

type encoder struct {
	plaintext string
	bytes     []byte
}

type decoder struct {
	base64 string
}

func NewEncoder(plaintext *string) *encoder {
	if plaintext == nil {
		plaintext = new(string)
		*plaintext = ""
	}
	return &encoder{
		plaintext: *plaintext,
		bytes:     []byte(*plaintext),
	}
}

func (e *encoder) Plaintext() string {
	return e.plaintext
}

func (e *encoder) Base64() string {
	return base64.RawURLEncoding.Strict().EncodeToString(e.bytes)
}

func (e *encoder) Hash(alg crypto.Hash) []byte {
	h := alg.New()
	h.Write([]byte(e.plaintext))
	return h.Sum(nil)
}

func (e *encoder) Bytes() []byte {
	return e.bytes
}

func (e *encoder) Random(length int) *encoder {
	bytes := make([]byte, length)
	rand.Read(bytes)
	e.bytes = bytes
	e.plaintext = string(e.bytes)
	return e
}

func (e *encoder) ReadBytes(bytes []byte) *encoder {
	e.bytes = bytes
	e.plaintext = string(e.bytes)
	return e
}

func (e *encoder) Hex(alg crypto.Hash) (string) {
	if (alg == 0) {
		return hex.EncodeToString(e.Bytes())
	}
	return hex.EncodeToString(e.Hash(alg))
}

func NewDecoder(b64 string) *decoder {
	return &decoder{
		base64: b64,
	}
}

func (d *decoder) Bytes() ([]byte, error) {
	b, err := base64.RawURLEncoding.Strict().DecodeString(d.base64)
	if err == nil {
		return b, nil
	}
	return base64.StdEncoding.Strict().DecodeString(d.base64)
}

