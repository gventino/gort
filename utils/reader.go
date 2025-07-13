package utils

import (
	"bufio"
)

func ReadLine(r *bufio.Reader) ([]byte, error) {
	var b byte
	var err error
	var word []byte
	for b != '\n' {
		b, err = r.ReadByte()
		if err != nil {
			return word, err
		}
		word = append(word, b)
	}
	return word, nil
}

func ReadPasswords(r *bufio.Reader, max int) ([][]byte, error) {
	var buffer [][]byte
	var err error

	for len(buffer) < max {
		var line []byte
		line, err = ReadLine(r)

		if len(line) > 0 {
			buffer = append(buffer, line)
		}

		if err != nil {
			return buffer, err
		}
	}

	return buffer, nil
}
