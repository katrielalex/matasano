package matasano

import "bytes"
import "crypto/aes"
import "fmt"
import "log"
import "strings"
import "testing"

// import "unicode/utf8"

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
	for i := 0; i < 50; i++ {
		var wasEcb bool
		oracle := func(plain []byte) []byte {
			tmp, ciphertext := encryptionOracle(plain)
			wasEcb = tmp
			return ciphertext
		}
		// gotta hardcode 16 here since otherwise we'll call the oracle repeatedly to guess the
		// block size, but it does different things each time
		guessEcb := isEcbBs(oracle, 16)
		if guessEcb != wasEcb {
			t.Error("Failed to identify ECB")
		}
	}
}

func Test_2_12(t *testing.T) {
	// preliminaries
	bs := getBlockSize(ecbSecretKeyOracle)
	if !isEcb(ecbSecretKeyOracle) {
		t.Error("It definitely is ECB")
	}

	// for each block from i=0 to n_blocks...
	var decrypted bytes.Buffer
	unmodifiedCiphertext := ecbSecretKeyOracle([]byte(""))
	log.Print(len(unmodifiedCiphertext))
	defer func() { log.Printf("%d %q", len(decrypted.String()), decrypted.String()) }()
	for i := 0; i < len(unmodifiedCiphertext)/bs; i++ {
		// decrypt that block
		possibles := make(map[string]string)
		var plaintextBlock bytes.Buffer
		for shift := 1; shift <= bs; shift++ {
			// Plaintext so far, or AAAAs if we don't have any
			var prefix []byte
			if decrypted.Len() > 0 {
				prefix = decrypted.Bytes()[decrypted.Len()-bs+shift:]
			} else {
				prefix = []byte(strings.Repeat("A", bs-shift))
			}

			// Try one byte at a time to fill out a block
			for j := 0; j < 256; j++ {
				plaintext := append(prefix, plaintextBlock.Bytes()...)
				plaintext = append(plaintext, byte(j))
				ciphertext := ecbSecretKeyOracle(plaintext)[:bs]
				possibles[string(ciphertext)] = string(plaintext)
			}
			target := ecbSecretKeyOracle(prefix)[i*bs : (i+1)*bs]
			// log.Printf("i (%d), shift (%d), target len %d %q\n%q", i, shift, len(target), prefix, target)
			c := possibles[string(target)]
			if c == "" {
				panic("Didn't find ciphertext block in dictionary")
			} else {
				plaintextBlock.Write([]byte(c[bs-1:]))
			}
		}
		decrypted.Write(plaintextBlock.Bytes())
	}
}
