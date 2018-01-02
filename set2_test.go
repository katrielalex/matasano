package matasano

import "bytes"
import "crypto/aes"
import "fmt"

// import "log"
import "strings"
import "testing"

func Test_2_9(t *testing.T) {
	got := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x04\x04\x04\x04")
	padded := pkcs7(got, 20)
	if !bytes.Equal(want, padded) {
		t.Errorf("You got PKCS7 wrong: expected %q but got %q",
			want,
			padded)
	}
}

func Test_2_10(t *testing.T) {
	b, err := aes.NewCipher([]byte("YELLOW SUBMARINE"))
	check(err)
	iv := []byte(strings.Repeat("\x00", 16))
	ciphertext := readB64File("data/10.txt")
	plaintext := fmt.Sprintf("%q", cbcDecrypt(b, iv, ciphertext))
	if !strings.Contains(plaintext, "Cause why the freaks are jockin' like Crazy Glue") {
		t.Error("Failed CBC decryption")
	}
}

func Test_cbc_roundtrip(t *testing.T) {
	b, err := aes.NewCipher([]byte("YELLOW SUBMARINE"))
	check(err)
	plaintext := []byte(strings.Repeat("ORANGE SUBMARINE", 100))
	iv := []byte(strings.Repeat("\x00", 16))
	ciphertext := cbcEncrypt(b, iv, plaintext)
	if !bytes.Equal(cbcDecrypt(b, iv, ciphertext), plaintext) {
		t.Error("Failed CBC roundtrip")
	}
}

func Test_random_does_not_repeat(t *testing.T) {
	const bs int = 16
	if bytes.Equal(randomAESKey(bs), randomAESKey(bs)) {
		t.Error("You have RNG problems")
	}
}

func Test_encryptionOracle(t *testing.T) {
	plain := []byte("sup world")
	_, ciphertext := encryptionOracle(plain)
	if bytes.Equal(plain, ciphertext) {
		t.Error("You done goofed")
	}
	// log.Printf("\n%d %q\n%d %q", len(plain), plain, len(ciphertext), ciphertext)
}

func Test_2_11(t *testing.T) {
	const bs int = 16
	plain := []byte(strings.Repeat("\x00", 1024))
	for i := 0; i < 50; i++ {
		wasEcb, ciphertext := encryptionOracle(plain)
		guessEcb := hasRepeatedBlocks(ciphertext, bs)
		if guessEcb != wasEcb {
			t.Error("Failed to identify ECB")
		}
	}
}
