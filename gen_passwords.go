package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

func GeneratePasswords(path string, size, qty int) {
	// CPU PROFILING
	var err error
	fCPU, err := os.Create("profs/password_cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(fCPU)
	defer pprof.StopCPUProfile()

	fmt.Printf("Starting the generation of %d passwords\n\tUsing: %d workers\n", qty, runtime.NumCPU())

	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	passwordsChan := make(chan string, 100)
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobsPerWorker := int(qty / numWorkers)

	// producer
	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobsPerWorker {
				passwordsChan <- GeneratePassword(size)
			}
		}()
	}

	// closes the chan
	go func() {
		wg.Wait()
		close(passwordsChan)
	}()

	// consumer
	for password := range passwordsChan {
		_, err := fmt.Fprintf(f, "%s\n", password)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Passwords generation done!")

	// MEM PROFILING
	fMEM, err := os.Create("profs/password_mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMEM.Close()
	if err := pprof.WriteHeapProfile(fMEM); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func GeneratePassword(size int) string {
	buffer := make([]byte, size)
	for i := range size {
		buffer[i] = byte(CHARSET[rand.IntN(len(CHARSET))])
	}
	return string(buffer)
}
