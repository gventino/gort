package utils

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

func BytesToOption(bytes []byte) (int, error) {
	if len(bytes) == 1 {
		b := bytes[0]
		// yeah baby, ascii codes
		if b > 48 && b < 57 {
			return int(b) - 48, nil
		}
		return -1, fmt.Errorf("error: cannot parse bytes to option")
	}
	return -1, fmt.Errorf("error: cannot parse bytes to option")
}

func Chain(b []byte, h hash.Hash, n int, table map[string]string) {
	var hashcode []byte
	plaintext := string(b)
	for range n {
		h.Reset()
		h.Write([]byte(plaintext))
		hashcode = h.Sum(nil)
		plaintext = hex.EncodeToString(hashcode)
	}
	table[hex.EncodeToString(hashcode)] = string(b)
}

func ResetReader(file *os.File, r *bufio.Reader) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	r.Reset(file)
}
