package matasano

import "bufio"
import "bytes"
import "log"
import "os"
import "strings"
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
	_, plaintext := anglify(x)
	log.Print(plaintext)
}

func Test_1_4(t *testing.T) {
	f, err := os.Open("data/4.txt")
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	score, plaintext := 0, ""
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\n")
		guessScore, guessPlaintext := anglify(line)
		if guessScore > score {
			score, plaintext = guessScore, guessPlaintext
		}
	}
	log.Print(plaintext)
}

func Test_1_5(t *testing.T) {
	got := `Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`
	want := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
	// scanner := bufio.NewScanner(strings.NewReader(s))
	// for scanner.Scan() {
	// 	log.Print(scanner.Text())
	// 	log.Print(hexOfBytes(xorKey(scanner.Bytes(), "ICE")))
	// }
	if want != hexOfBytes(xorKey([]byte(got), "ICE")) {
		t.Errorf("Didn't get expected output from repeated XOR")
	}
}
