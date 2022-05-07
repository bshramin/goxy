package hashid

import (
	"testing"
)

var salt = "qtyq68eqeqwy"
var chars = "abcdefghijklmnopqrstuvwxyz1234567890"
var minLen = 6
func TestEncode(t *testing.T) {
	hasher := getHasher(t)
	id := int64(80175)
	code, err := hasher.Encode(id)
	if err != nil {
		t.Fatal("encoding failed")
	}
	if code != "3drnmd" {
		t.Fatal("encoding result is wrong")
	}
}

func TestDecode(t *testing.T) {
	hasher := getHasher(t)
	code := "3drnmd"
	id, err := hasher.Decode(code)
	if err != nil {
		t.Fatal("decoding failed")
	}
	if id != 80175 {
		t.Fatal("decoding result is wrong")
	}
}

func getHasher(t *testing.T) *hsh {
	hasher, err := New(salt, chars, minLen)
	if err != nil {
		t.Fatal("hasher init failed")
	}
	return hasher
}