// +build gofuzz

package reader

import "bytes"

func Fuzz(data []byte) int {
	_, err := New(bytes.NewBuffer(data)).One()
	if err != nil {
		return 0
	}

	return 1
}
