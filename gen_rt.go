package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

func GenerateTable(passwordsPath, binPath string) {
	// CPU PROFILING
	var err error
	fCPU, err := os.Create("profs/rt_cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(fCPU)
	defer pprof.StopCPUProfile()

	fmt.Printf("\nBUILDING RAINBOW TABLE\n")

	// build process
	file, err := os.OpenFile(passwordsPath, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := bufio.NewReader(file)

	table := make(RainbowTable, NUM_PASSWORDS)
	FillTable(r, table)
	fmt.Println("Rainbow table finished!")

	fmt.Println("Serializing rainbow table...")
	bin, err := os.Create(binPath)
	if err != nil {
		panic(err)
	}
	defer bin.Close()

	encoder := gob.NewEncoder(bin)
	if err = encoder.Encode(table); err != nil {
		panic(err)
	}
	fmt.Printf("Rainbow table serialized!\n")

	// MEM PROFILING
	fMEM, err := os.Create("profs/rt_mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMEM.Close()
	if err := pprof.WriteHeapProfile(fMEM); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func FillTable(r *bufio.Reader, table RainbowTable) RainbowTable {
	processedCount := 0
	hasher := sha512.New()

	// buffer alloc
	buffer := make([][]byte, PASSWORDS_BUFF_SIZE)
	for i := range len(buffer) {
		buffer[i] = make([]byte, PASSWORD_LENGTH)
	}

	var err error
	buffer, err = ReadPasswords(r, buffer, PASSWORDS_BUFF_SIZE)

	// bufferized file reading
	var start string
	var end []byte
	for err != io.EOF {
		for i := range len(buffer) {
			start = string(buffer[i])
			end = Chain(buffer[i], hasher)
			table[string(end)] = string(start)
			processedCount++
		}

		if processedCount%50000 == 0 && processedCount > 0 {
			fmt.Printf("Processed %d passwords\n", processedCount)
		}

		buffer, err = ReadPasswords(r, buffer, PASSWORDS_BUFF_SIZE)
		if err != nil {
			break
		}
	}
	fmt.Printf("Table Completed: %d chains generated\n", len(table))
	return table
}

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

func ReadPasswords(r *bufio.Reader, buffer [][]byte, max int) ([][]byte, error) {
	var err error

	for i := range max {
		var line []byte
		line, err = ReadLine(r)

		if len(line) > 0 {
			buffer[i] = line
		}

		if err != nil {
			return buffer, err
		}
	}

	return buffer, nil
}
