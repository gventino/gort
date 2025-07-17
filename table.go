package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

type RainbowTable map[string]string

func Hash(plaintext []byte, hasher hash.Hash) []byte {
	hasher.Reset()
	hasher.Write(plaintext)
	return hasher.Sum(nil)
}

func Reduce(hashcode []byte, pos int) []byte {
	pw := make([]byte, PASSWORD_LENGTH)
	for i := range PASSWORD_LENGTH {
		idx := (int(hashcode[i%len(hashcode)]) + pos + i) % len(CHARSET)
		pw[i] = CHARSET[idx]
	}
	return pw
}

func Chain(startPassword []byte, hasher hash.Hash) []byte {
	currentPassword := startPassword
	var hashVal []byte
	for i := range CHAIN_SIZE {
		hashVal = Hash(currentPassword, hasher)
		currentPassword = Reduce(hashVal, i)
	}
	return currentPassword
}

func Lookup(targetHash string, table RainbowTable) ([]byte, bool) {
	var currentHash, tempPassword, hashVal, regenPassword, currentRegenHash []byte
	hasher := sha512.New()
	for pos := CHAIN_SIZE - 1; pos >= 0; pos-- {
		var err error
		currentHash, err = hex.DecodeString(targetHash)
		if err != nil {
			return []byte(""), false
		}

		tempPassword = Reduce(currentHash, pos)
		for i := pos + 1; i < CHAIN_SIZE; i++ {
			hashVal = Hash(tempPassword, hasher)
			tempPassword = Reduce(hashVal, i)
		}

		if startPassword, found := table[string(tempPassword)]; found {
			regenPassword = []byte(startPassword)
			for i := range CHAIN_SIZE {
				currentRegenHash = Hash(regenPassword, hasher)
				currentRegenHashHex := fmt.Sprintf("%x", currentRegenHash)
				if currentRegenHashHex == targetHash {
					return regenPassword, true
				}
				regenPassword = Reduce(currentRegenHash, i)
			}
		}
	}

	return []byte(""), false
}
