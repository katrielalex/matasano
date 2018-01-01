from utils import Challenge, solves, \
    hex_to_bytes, bytes_to_b64, hex_to_b64, xor


@solves(1, 1)
def problem_1(
        s=b"49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d",
        t=b"SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t",
):
    assert hex_to_b64(s) == t


@solves(1, 2)
def problem_2():
    s = b"1c0111001f010100061a024b53535009181c"
    t = b"686974207468652062756c6c277320657965"
    u = b"746865206b696420646f6e277420706c6179"

    assert xor(hex_to_bytes(s), hex_to_bytes(t)) == hex_to_bytes(u)
