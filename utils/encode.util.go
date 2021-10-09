package utils

import "encoding/base64"

var stdEncode = base64.RawStdEncoding.Strict()
var urlEncode = base64.RawURLEncoding.Strict()

func Base64(dst, src []byte) {
	stdEncode.Encode(dst, src)
}

func Base64URL(dst, src []byte) {
	urlEncode.Encode(dst, src)
}

func Base64String(b []byte) string {
	return stdEncode.EncodeToString(b)
}

func Base64URLString(b []byte) string {
	return urlEncode.EncodeToString(b)
}

func Base64StringDecode(s string) ([]byte, error) {
	b, err := stdEncode.DecodeString(s)
	return b, err
}

func Base64URLStringDecode(s string) ([]byte, error) {
	b, err := urlEncode.DecodeString(s)
	return b, err
}

func Base64URLStringBulkDecode(s ...string) ([][]byte, error) {
	bs := make([][]byte, 0)
	for _, value := range s {
		b, err := Base64URLStringDecode(value)
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

func Base64StringDecodeToString(s string) (string, error) {
	b, err := Base64StringDecode(s)
	return string(b), err
}

func Base64URLStringDecodeToString(s string) (string, error) {
	b, err := Base64URLStringDecode(s)
	return string(b), err
}
