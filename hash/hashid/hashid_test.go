package hashid

import (
	"fmt"
	"testing"
)

var (
	salt   = "qtyq68eqeqwy"
	chars  = "abcdefghijklmnopqrstuvwxyz1234567890"
	minLen = 6
)

func TestEncode(t *testing.T) {
	t.Parallel()
	hasher, err := getHasher()
	if err != nil {
		t.Fatal(err.Error())
	}
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
	t.Parallel()
	hasher, err := getHasher()
	if err != nil {
		t.Fatal(err.Error())
	}
	code := "3drnmd"
	id, err := hasher.Decode(code)
	if err != nil {
		t.Fatal("decoding failed")
	}
	if id != 80175 {
		t.Fatal("decoding result is wrong")
	}
}

func getHasher() (*hsh, error) {
	hasher, err := New(salt, chars, minLen)
	if err != nil {
		return nil, fmt.Errorf("hasher init failed")
	}
	return hasher, nil
}
