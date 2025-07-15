package utils

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#%^&*?")
const passwordLength = 5

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

func Reduce(hashBytes []byte, pos int) string {
	pw := make([]rune, passwordLength)
	for i := 0; i < passwordLength; i++ {
		idx := (int(hashBytes[i%len(hashBytes)]) + pos + i) % len(charset)
		pw[i] = charset[idx]
	}
	return string(pw)
}

func Chain(start []byte, h hash.Hash, n int, table map[string]string) {
	plaintext := string(start)
	var hashcode []byte
	for i := 0; i < n; i++ {
		h.Reset()
		h.Write([]byte(plaintext))
		hashcode = h.Sum(nil)
		plaintext = Reduce(hashcode, i)
	}
	table[plaintext] = string(start)
}

func ChainReturnEnd(start []byte, h hash.Hash, n int) string {
	plaintext := string(start)
	var hashcode []byte
	for i := 0; i < n; i++ {
		h.Reset()
		h.Write([]byte(plaintext))
		hashcode = h.Sum(nil)
		plaintext = Reduce(hashcode, i)
	}
	return plaintext
}

func ResetReader(file *os.File, r *bufio.Reader) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	r.Reset(file)
}

type ChainEntry struct {
	Start string
	End   string
}

func LookupInRainbowTable(targetHash string, table []ChainEntry, h hash.Hash, chainLen int) (string, bool) {
	for pos := chainLen - 1; pos >= 0; pos-- {
		temp := targetHash
		for i := pos; i < chainLen; i++ {
			// Reduction
			reduced := Reduce(decodeHex(temp), i)
			// Hash
			h.Reset()
			h.Write([]byte(reduced))
			hashcode := h.Sum(nil)
			temp = encodeHex(hashcode)
		}
		// check if temp matches any chain end
		for _, entry := range table {
			if entry.End == temp {
				// rebuild chain from beginning
				plaintext := entry.Start
				for i := 0; i < chainLen; i++ {
					h.Reset()
					h.Write([]byte(plaintext))
					hashcode := h.Sum(nil)
					if encodeHex(hashcode) == targetHash {
						return plaintext, true
					}
					plaintext = Reduce(hashcode, i)
				}
			}
		}
	}
	return "", false
}

func decodeHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func encodeHex(b []byte) string {
	return hex.EncodeToString(b)
}

func EncodeHex(b []byte) string {
	return hex.EncodeToString(b)
}

func GetSHA256() hash.Hash {
	return sha256.New()
}

func GetSHA512() hash.Hash {
	return sha512.New()
}

func GetSHA384() hash.Hash {
	return sha512.New384()
}

func GetSHA3_256() hash.Hash {
	return sha3.New256()
}

func GetSHA3_512() hash.Hash {
	return sha3.New512()
}

func GetMD5() hash.Hash {
	return md5.New()
}
