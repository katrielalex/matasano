package matasano

import b64 "encoding/base64"
import hex "encoding/hex"
import "log"

// decode a hex string into a byte array
func bytes_of_hex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

// encode a byte array into b64
func b64_of_bytes(b []byte) string {
	return b64.StdEncoding.EncodeToString(b)
}

// convert hex to b64 via bytes
func b64_of_hex(s string) string {
	return b64_of_bytes(bytes_of_hex(s))
}

// xor two byte arrays
func xor(a, b []byte) []byte {
	n := len(a)
	if n > len(b) { n = len(b) }
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}
