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

func main() {
	var err error
	// ignore, profiling
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	var file, bin *os.File
	hashes := []hash.Hash{
		sha256.New(),
		sha512.New(),
		sha512.New384(),
		sha3.New256(),
		sha3.New512(),
		sha3.New384(),
		md5.New(),
	}
	rt := make([]map[string]string, len(hashes))
	var table map[string]string

	// reading passwords list
	file, err = os.OpenFile("lists/jwt.secrets.list", os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// building rt
	fmt.Println("CONSTRUINDO RAINBOW TABLE!")
	r := bufio.NewReader(file)
	for i := range len(hashes) {
		table = GenerateTable(r, hashes[i])
		rt[i] = table
		utils.ResetReader(file, r)
	}
	fmt.Println("RAINBOW TABLE CONSTRUÍDA")

	// saving the rt
	fmt.Println("SERIALIZANDO RAINBOW TABLE")
	bin, err = os.Create("rainbow_tables/jwt.secrets.bin")
	if err != nil {
		panic(err)
	}
	defer bin.Close()

	encoder := gob.NewEncoder(bin)
	if err = encoder.Encode(rt); err != nil {
		panic(err)
	}
	fmt.Println("RT SERIALIZADA COM SUCESSO!")

	// ignore, profiling
	fMEM, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMEM.Close()

	if err := pprof.WriteHeapProfile(fMEM); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func GenerateTable(r *bufio.Reader, h hash.Hash) map[string]string {
	table := make(map[string]string)
	var buffer [][]byte
	var err error

	buffer, err = utils.ReadPasswords(r, cmd.PASSWORDS_BUFF_SIZE)
	for err != io.EOF {
		for i := range len(buffer) {
			utils.Chain(buffer[i], h, cmd.CHAIN_SIZE, table)
		}
		buffer, err = utils.ReadPasswords(r, cmd.PASSWORDS_BUFF_SIZE)
		if err != nil {
			break
		}
	}
	return table
}
