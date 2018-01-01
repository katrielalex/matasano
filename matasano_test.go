package matasano

import "bytes"
import "fmt"

// import "log"
import "testing"

func Test_1_1(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d", "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"},
		{"", ""},
	}
	for _, c := range cases {
		got := b64OfHex(c.in)
		if got != c.want {
			t.Errorf("b64OfHex(%q) == %q, want %q",
				c.in,
				got,
				c.want)
		}
	}
}

func Test_1_2(t *testing.T) {
	x := bytesOfHex("1c0111001f010100061a024b53535009181c")
	y := bytesOfHex("686974207468652062756c6c277320657965")
	z := bytesOfHex("746865206b696420646f6e277420706c6179")

	if !bytes.Equal(xor(x, y), z) {
		t.Errorf("xoring gave the wrong answer")
	}
}

func Test_1_3(t *testing.T) {
	x := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	xB := bytesOfHex(x)

	plaintext := xB
	score := englishnessCount(xB)
	for rune := range englishFreqs {
		guess := xorcs(xB, rune)
		guessScore := englishnessCount(guess)
		if guessScore > score {
			plaintext, score = guess, guessScore
		}
	}
	if fmt.Sprintf("%s", plaintext) == "" {
		t.Errorf("wrongo")
	}
	// log.Printf("%s", plaintext)
}
