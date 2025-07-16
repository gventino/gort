package main

import (
	"crypto/sha512"
	"fmt"
	"hash"
)

type RainbowTable map[string]string

func Hash(plaintext []byte, hasher hash.Hash) []byte {
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
	for i := range CHAIN_SIZE {
		hashVal := Hash(currentPassword, hasher)
		currentPassword = Reduce(hashVal, i)
	}
	return currentPassword
}

func Lookup(targetHash string, table RainbowTable) ([]byte, bool) {
	fmt.Printf("Searching password for hashcode: %s\n", targetHash)
	var currentHash, tempPassword, hashVal, regenPassword, currentRegenHash []byte
	hasher := sha512.New()
	for pos := CHAIN_SIZE - 1; pos >= 0; pos-- {
		currentHash = []byte(targetHash)

		tempPassword = Reduce(currentHash, pos)
		for i := pos + 1; i < CHAIN_SIZE; i++ {
			hashVal = Hash(tempPassword, hasher)
			tempPassword = Reduce(hashVal, i)
		}

		if startPassword, found := table[string(tempPassword)]; found {
			fmt.Println("Possible Match! Recreating chain from:\n\t", startPassword)
			regenPassword = []byte(startPassword)
			for i := range CHAIN_SIZE {
				currentRegenHash = Hash(regenPassword, hasher)
				if string(currentRegenHash) == targetHash {
					return regenPassword, true
				}
				regenPassword = Reduce(currentRegenHash, i)
			}
			fmt.Println("Fake News!!! Target hashcode not found in Chain")
		}
	}

	return []byte(""), false
}
