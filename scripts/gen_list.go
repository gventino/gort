package main

import (
	"main/utils"
	"os"
	"strconv"
	"sync"
)

type seed struct {
	value uint64
}

func (s *seed) Uint64() uint64 {
	return s.value
}

func main() {
	args := os.Args[1:]

	size, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}

	qty, err := strconv.Atoi(args[1])
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	numGoroutines := 16
	wg.Add(numGoroutines)
	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range int(qty / numGoroutines) {
				utils.GeneratePassword(size)
			}
		}()
	}

	wg.Wait()
}
