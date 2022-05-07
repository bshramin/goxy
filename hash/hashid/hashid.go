package hashid

import (
	"github.com/speps/go-hashids"
)

type hsh struct {
	hashids *hashids.HashID
}

func New(salt, chars string, len int) (*hsh, error) {
	hidsData := hashids.NewData()
	hidsData.Alphabet = chars
	hidsData.Salt = salt
	hidsData.MinLength = len
	h, err := hashids.NewWithData(hidsData)
	if err != nil {
		return nil, err
	}
	res := &hsh{
		hashids: h,
	}
	return res, nil
}

func (h *hsh) Encode(id int64) (string, error) {
	numbers := []int{int(id)}
	res, err := h.hashids.Encode(numbers)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (h *hsh) Decode(code string) (int64, error) {
	ii, err := h.hashids.DecodeInt64WithError(code)
	if err != nil {
		return 0, err
	}
	res := ii[0]
	return res, nil
}
