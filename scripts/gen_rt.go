package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"hash"
	"io"
	"log"
	"main/cmd"
	"main/utils"
	"os"
	"runtime/pprof"
)

type ChainEntry struct {
	Start string
	End   string
}

func main() {
	var err error
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	hashes := []hash.Hash{
		sha256.New(),
		sha512.New(),
		sha512.New384(),
		sha3.New256(),
		sha3.New512(),
		sha3.New384(),
		md5.New(),
	}
	rt := make([][]ChainEntry, len(hashes))

	fmt.Printf("\nCONSTRUINDO RAINBOW TABLE\n")
	fmt.Printf("Chain length: %d\n", cmd.CHAIN_SIZE)

	for i := range len(hashes) {
		fmt.Printf("Building table %d/%d (Hash function: %T)\n", i+1, len(hashes), hashes[i])

		var allChains []ChainEntry

		file, err := os.OpenFile("lists/jwt.secrets.list", os.O_RDONLY, 0777)
		if err != nil {
			panic(err)
		}

		r := bufio.NewReader(file)
		originalChains := GenerateTable(r, hashes[i], "original passwords")
		allChains = append(allChains, originalChains...)
		file.Close()

		rt[i] = allChains
		fmt.Printf("  Total chains for this hash function: %d\n\n", len(allChains))
	}

	fmt.Println("RAINBOW TABLE CONSTRUÍDA!")

	totalChains := len(rt[0])
	totalPositions := totalChains * cmd.CHAIN_SIZE
	passwordSpace := 1680700000 // 70^5
	coverage := float64(totalPositions) / float64(passwordSpace) * 100.0
	fmt.Printf("- Actual coverage: %.4f%%\n\n", coverage)

	fmt.Println("SERIALIZANDO RAINBOW TABLE")
	bin, err := os.Create("rainbow_tables/jwt.secrets.bin")
	if err != nil {
		panic(err)
	}
	defer bin.Close()

	encoder := gob.NewEncoder(bin)
	if err = encoder.Encode(rt); err != nil {
		panic(err)
	}

	fmt.Printf("RT SERIALIZADA COM SUCESSO!\n")
	fmt.Printf("File size: rainbow_tables/jwt.secrets.bin\n")

	fMEM, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMEM.Close()

	if err := pprof.WriteHeapProfile(fMEM); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func GenerateTable(r *bufio.Reader, h hash.Hash, source string) []ChainEntry {
	table := make([]ChainEntry, 0)
	var buffer [][]byte
	var err error
	var processedCount int
	const PASSWORDS_BUFF_SIZE = 10000

	fmt.Printf("  Processing %s...\n", source)
	buffer, err = utils.ReadPasswords(r, PASSWORDS_BUFF_SIZE)
	for err != io.EOF {
		for i := range len(buffer) {
			start := string(buffer[i])
			end := utils.ChainReturnEnd(buffer[i], h, cmd.CHAIN_SIZE)
			table = append(table, ChainEntry{Start: start, End: end})
			processedCount++
		}

		if processedCount%50000 == 0 && processedCount > 0 {
			fmt.Printf("    Processed %d passwords from %s\n",
				processedCount, source)
		}

		buffer, err = utils.ReadPasswords(r, PASSWORDS_BUFF_SIZE)
		if err != nil {
			break
		}
	}

	fmt.Printf("  Completed %s: %d chains generated\n", source, len(table))
	return table
}
