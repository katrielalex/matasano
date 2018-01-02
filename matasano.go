package matasano

import "bufio"
import "bytes"
import "crypto/cipher"
import "errors"
import "fmt"
import "log"
import "math"
import "os"
import "strings"
import "unicode/utf8"
import b64 "encoding/base64"
import hex "encoding/hex"

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// decode a hex string into a byte array
func bytesOfHex(s string) []byte {
	decoded, err := hex.DecodeString(s)
	check(err)
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

// decode a b64 string into bytes
func bytesOfB64(s string) []byte {
	t, err := b64.StdEncoding.DecodeString(s)
	check(err)
	return t
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

func xorKey(a []byte, key string) []byte {
	dst := make([]byte, len(a))
	for i, b := range a {
		dst[i] = b ^ key[i%len(key)]
	}
	return dst
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
			if 'A' <= rune && rune <= 'Z' {
				total += int(englishFreqs[rune+26] * 5)
			} else {
				total += int(englishFreqs[rune] * 10)
			}
		} else {
			if strings.ContainsRune("' \n", rune) {
				// pass
			} else {
				total -= SymbolPenalty
			}
		}
	}
	return total
}

// xor with the character making it the most english-y
func anglifyB(xB []byte) (rune, int, string) {
	plaintext := xB
	score := englishnessCount(xB)
	key := '0'
	for rune := ' '; rune <= '~'; rune++ {
		guessPlaintext := xorcs(xB, rune)
		guessScore := englishnessCount(guessPlaintext)
		// log.Printf("%d, %d, %q", score, guessScore, guessPlaintext)
		if guessScore > score {
			score, plaintext = guessScore, guessPlaintext
			key = rune
		}
	}
	return key, score, fmt.Sprintf("%s", plaintext)
}

func anglify(x string) (rune, int, string) {
	xB := bytesOfHex(x)
	return anglifyB(xB)
}

func hamming(a, b []byte) (int, error) {
	if len(a) != len(b) {
		return 0, errors.New("Can't compute Hamming distance for strings of different length")
	}

	distance := 0
	// bytes
	for i := 0; i < len(a); i++ {
		// bits
		b1, b2 := a[i], b[i]
		for j := 0; j < 8; j++ {
			mask := byte(1 << uint(j))
			if (b1 & mask) != (b2 & mask) {
				distance++
			}
		}
	}
	return distance, nil
}

func hammingS(a, b string) (int, error) {
	return hamming([]byte(a), []byte(b))
}

func readB64File(path string) []byte {
	// I'm sure there's a better way to read a b64 file...
	f, err := os.Open(path)
	check(err)
	defer func() { check(f.Close()) }()

	scanner := bufio.NewScanner(f)
	var b bytes.Buffer
	for scanner.Scan() {
		b.Write(bytesOfB64(strings.TrimSuffix(scanner.Text(), "\n")))
	}
	s := b.Bytes()
	return s
}

func transposeInplace(a [][]byte) [][]byte {
	n := len(a)
	if n == 0 {
		panic("WTF are you doing with zero-length arrays?!")
	}
	m := len(a[0])

	dst := make([][]byte, m)
	for j := 0; j < m; j++ {
		dst[j] = make([]byte, n)
		for i := 0; i < n; i++ {
			dst[j][i] = a[i][j]
		}
	}
	return dst
}

func ecbDecrypt(block cipher.Block, ciphertext []byte) []byte {
	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		panic("Need a multiple of the blocksize")
	}
	var plain bytes.Buffer
	for i := 0; i < len(ciphertext)/bs; i++ {
		plaintext := make([]byte, bs)
		block.Decrypt(plaintext, ciphertext[i*bs:(i+1)*bs])
		plain.Write(plaintext)
	}
	return plain.Bytes()
}

func hasRepeatedBlocks(ciphertext []byte, bs int) bool {
	if len(ciphertext)%bs > 0 {
		panic("Need a multiple of the blocksize")
	}
	m := make(map[string]bool)
	for i := 0; i < len(ciphertext)/bs; i++ {
		m[string(ciphertext[i*bs:(i+1)*bs])] = false
	}
	return len(m) < len(ciphertext)/bs
}
