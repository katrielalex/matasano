package matasano

import "bytes"
import "crypto/aes"
import "fmt"

// import "log"
import "strings"
import "testing"

func Test_2_1(t *testing.T) {
	got := []byte("YELLOW SUBMARINE")
	want := []byte("YELLOW SUBMARINE\x04\x04\x04\x04")
	padded := pkcs7(got, 20)
	if !bytes.Equal(want, padded) {
		t.Errorf("You got PKCS7 wrong: expected %q but got %q",
			want,
			padded)
	}
}

func Test_2_2(t *testing.T) {
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
	plaintext := []byte("ORANGE SUBMARINE")
	iv := []byte(strings.Repeat("\x00", 16))
	ciphertext := cbcEncrypt(b, iv, plaintext)
	if !bytes.Equal(cbcDecrypt(b, iv, ciphertext), plaintext) {
		t.Error("Failed CBC roundtrip")
	}
}
