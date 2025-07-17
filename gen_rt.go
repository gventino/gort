package main

import (
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

func GenerateTable(_ string, binPath string) {
	// CPU PROFILING
	var err error
	fCPU, err := os.Create("profs/rt_cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(fCPU)
	defer pprof.StopCPUProfile()

	fmt.Printf("\nBUILDING RAINBOW TABLE\n")

	table := make(RainbowTable, NUM_PASSWORDS)
	FillTable(table)
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

func FillTable(table RainbowTable) RainbowTable {
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	jobs := make(chan struct{}, 1000)
	result := make(chan struct{ end, start string }, 1000)
	used := sync.Map{}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			hasher := sha512.New()
			for range jobs {
				var pw string
				for {
					pw = GeneratePassword(PASSWORD_LENGTH)
					if _, exists := used.LoadOrStore(pw, struct{}{}); !exists {
						break
					}
				}
				end := Chain([]byte(pw), hasher)
				result <- struct{ end, start string }{string(end), pw}
			}
			wg.Done()
		}()
	}

	go func() {
		for i := 0; i < NUM_PASSWORDS; i++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}()

	for i := 0; i < NUM_PASSWORDS; i++ {
		r := <-result
		table[r.end] = r.start
		if (i+1)%5000 == 0 || i+1 == NUM_PASSWORDS {
			percent := float64(i+1) / float64(NUM_PASSWORDS) * 100
			fmt.Printf("Processed %d/%d (%.2f%%) passwords\n", i+1, NUM_PASSWORDS, percent)
		}
	}
	wg.Wait()
	close(result)
	fmt.Printf("Table Completed: %d chains generated\n", len(table))
	return table
}
