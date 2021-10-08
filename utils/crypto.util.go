package utils

import (
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func HashString(plaintext string, alg crypto.Hash) []byte {
	h := alg.New()
	h.Write([]byte(plaintext))
	return h.Sum(nil)
}

func HashStringToHex(plaintext string, alg crypto.Hash) string {
	return fmt.Sprintf("%x", HashString(plaintext, alg))
}

func HashStringToBase64(plaintext string, alg crypto.Hash) string {
	return base64.RawStdEncoding.EncodeToString(HashString(plaintext, alg))
}

func HashStringToBase64URL(plaintext string, alg crypto.Hash) string {
	return base64.RawURLEncoding.EncodeToString(HashString(plaintext, alg))
}

func RandomBytes(n uint64) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
