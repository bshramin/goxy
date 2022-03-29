package goxy

import (
	"bytes"
	"encoding/gob"
)

func Encode(m interface{}) (string, error) {
	target := &bytes.Buffer{}
	err := gob.NewEncoder(target).Encode(m)
	if err != nil {
		return "", err
	}
	return target.String(), nil
}

func Decode[K any](input string) (K, error) {
	var res K
	buf := bytes.NewBufferString(input)
	err := gob.NewDecoder(buf).Decode(&res)
	return res, err
}
