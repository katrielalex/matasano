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
	_, _, plaintext := anglify(x)
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
		_, guessScore, guessPlaintext := anglify(line)
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

func Test_1_6a(t *testing.T) {
	x := "this is a test"
	y := "wokka wokka!!!"

	wanted := 37
	got, err := hammingS(x, y)
	check(err)
	if wanted != got {
		t.Errorf("Hamming distance calculator is buggy, expected %d but got %d", wanted, got)
	}
}

func Test_1_6b(t *testing.T) {
	// I'm sure there's a better way to read a b64 file...
	f, err := os.Open("data/6.txt")
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var b bytes.Buffer
	for scanner.Scan() {
		b.Write(bytesOfB64(strings.TrimSuffix(scanner.Text(), "\n")))
	}
	s := b.Bytes()

	// Find the keysize as the length with the smallest normalised edit
	// distance between pairs of blocks
	size, score := 2, 1000000.0
	for keysize := 4; keysize <= 40; keysize++ { // gets hung up on keysize=2/3
		const blocksToAverage int = 3
		guessScore := 0.0
		for i := 0; i < blocksToAverage; i++ {
			dist, err := hamming(s[i*keysize:(i+1)*keysize],
				s[(i+1)*keysize:(i+2)*keysize])
			check(err)
			// log.Printf("normalised distance %f between block %q and block %q", float64(dist) / float64(keysize), s[i*keysize:(i+1)*keysize], s[(i+1)*keysize:(i+2)*keysize])
			guessScore += float64(dist) / float64(keysize)
		}
		guessScore /= float64(blocksToAverage)

		// log.Printf("score (%f), guessScore (%f)", score, guessScore)
		if guessScore < score {
			score, size = guessScore, keysize
		}

	}

	// Pad if necessary (turns out it isn't)
	remainder := len(s) % size
	if remainder > 0 {
		s = append(s, []byte(strings.Repeat(" ", remainder))...)
		// log.Printf("Padding with extra %d bytes", remainder)
	}

	// Transpose blocks to get chunks xored with the same key
	chunks := make([][]byte, len(s)/size)
	for i := 0; i < len(s)/size; i++ {
		chunks[i] = s[i*size : (i+1)*size]
	}

	key := ""
	chunks = transposeInplace(chunks)
	for i := 0; i < len(chunks); i++ {
		keyC, _, plain := anglifyB(chunks[i])
		key += string(keyC)
		chunks[i] = []byte(plain)
	}
	chunks = transposeInplace(chunks)

	var plain bytes.Buffer
	for _, chunk := range chunks {
		plain.Write(chunk)
	}
	log.Print(key)
	// log.Print(plain.String())
}
