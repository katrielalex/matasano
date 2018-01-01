package matasano

import b64 "encoding/base64"
import hex "encoding/hex"
import "fmt"
import "log"
import "math"
import "unicode/utf8"

// decode a hex string into a byte array
func bytesOfHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}
	return decoded
}

// encode byte array to hex string
func hexOfBytes(b []byte) string {
	return hex.EncodeToString(b)
}

// encode a byte array into b64
func b64OfBytes(b []byte) string {
	return b64.StdEncoding.EncodeToString(b)
}

// convert hex to b64 via bytes
func b64OfHex(s string) string {
	return b64OfBytes(bytesOfHex(s))
}

// xor two byte arrays
func xor(a, b []byte) []byte {
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	dst := make([]byte, n)
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return dst
}

// xor byte array and char
func xorc(a []byte, b byte) []byte {
	dst := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		dst[i] = a[i] ^ b
	}
	return dst
}

// same as the above but with a rune
func xorcs(a []byte, rune rune) []byte {
	b := make([]byte, 1)
	if utf8.RuneLen(rune) > 1 {
		log.Fatal("wtf")
	}
	utf8.EncodeRune(b, rune)
	return xorc(a, b[0])
}

// copy-pasted from wikipedia
var englishFreqs = map[rune]float64{
	'a': 0.08167,
	'b': 0.01492,
	'c': 0.02782,
	'd': 0.04253,
	'e': 0.12702,
	'f': 0.02228,
	'g': 0.02015,
	'h': 0.06094,
	'i': 0.06966,
	'j': 0.00153,
	'k': 0.00772,
	'l': 0.04025,
	'm': 0.02406,
	'n': 0.06749,
	'o': 0.07507,
	'p': 0.01929,
	'q': 0.00095,
	'r': 0.05987,
	's': 0.06327,
	't': 0.09056,
	'u': 0.02758,
	'v': 0.00978,
	'w': 0.02360,
	'x': 0.00150,
	'y': 0.01974,
	'z': 0.0007,
}

// english-ness score by character freq
func englishnessL2(p []byte) float64 {
	n := len(p)
	nInv := 1.0 / float64(n)
	freqs := make(map[rune]float64)
	for _, rune := range fmt.Sprintf("%s", p) {
		freqs[rune] += nInv
	}

	// l2 difference between frequencies
	total := 0.0
	for rune, freq := range englishFreqs {
		diff := freq - freqs[rune]
		total += diff * diff
	}
	total = math.Sqrt(total)

	return -total
}

func isLetter(r rune) bool {
	return ('A' <= r && r <= 'Z') || ('a' <= r && r <= 'z')
}

// english-ness score by symbol counting: +frequency points for a letter, -1
// point for anything not a-zA-Z
func englishnessCount(p []byte) int {
	const SymbolPenalty int = 1
	total := 0
	for _, rune := range fmt.Sprintf("%s", p) {
		if isLetter(rune) {
			if rune > 'Z' {
				rune -= 26
			}
			total += int(englishFreqs[rune] * 100)
		} else {
			total -= SymbolPenalty
		}
	}
	return total
}
