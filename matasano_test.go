package matasano

import "bytes"
import "log"
import "testing"

func Test_1_1(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d", "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"},
		{"", ""},
	}
	for _, c := range cases {
		got := b64_of_hex(c.in)
		if got != c.want {
			t.Errorf("b64_of_hex(%q) == %q, want %q",
				c.in,
				got,
				c.want)
		}
	}
}

func Test_1_2(t *testing.T) {
	x := bytes_of_hex("1c0111001f010100061a024b53535009181c")
	y := bytes_of_hex("686974207468652062756c6c277320657965")
	z:= bytes_of_hex("746865206b696420646f6e277420706c6179")

	if !bytes.Equal(xor(x, y), z) {
		t.Errorf("xoring gave the wrong answer")
	}
}

func Test_1_3(t *testing.T) {
	x := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	x_b := bytes_of_hex(x)

	plaintext := x_b
	score := englishness_count(x_b)
	for rune, _ := range english_freqs {
		guess := xorc_s(x_b, rune)
		guess_score := englishness_count(guess)
		if guess_score > score {
			plaintext, score = guess, guess_score
		}
	}
	log.Printf("%s", plaintext)
}
