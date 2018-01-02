package matasano

import "bytes"
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
