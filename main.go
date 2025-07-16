package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

func main() {
	fmt.Print(ASCII)
	fmt.Print(MENU)

	// scanning first line
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	opt, _ := BytesToOption(scanner.Bytes())

	for opt != 4 {
		// cleaning stdout buffer with ansi scape code
		fmt.Print("\033[H\033[2J")
		switch opt {
		// capture the entries with []byte always, there are lighter than string.
		// and the https://pkg.go.dev/crypto (native go lib for crypto hashing) works by default with []byte;
		case 1:
			GeneratePasswords(PASSWORDS_PATH, PASSWORD_LENGTH, NUM_PASSWORDS)
		case 2:
			GenerateTable(PASSWORDS_PATH, BIN_PATH)
		case 3:
			DecryptHashcode()
		default:
			// continue asking for input
			fmt.Println("error: invalid option")
		}
		// this is only because we're cleaning the console with ansi scape chars
		time.Sleep(time.Second * 1)
		fmt.Print(ASCII)
		fmt.Print(MENU)
		scanner.Scan()
		opt, _ = BytesToOption(scanner.Bytes())
	}
	fmt.Println("CLOSING GoRt, THANKS FOR USING!")
	// this looks very cool, i like it pretty much
	time.Sleep(time.Millisecond * 500)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 500)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 500)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 500)
	fmt.Print("\033[H\033[2J")
}

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

func DecryptHashcode() {
	hashcode := "6f438c5ce2bf7073fcf6fcd9c2a2dfb5faa60b919dc87a62fc2d00eb08d025aca662a28e42c44ae8f29a7fb610920b2f7f3cee0346e80e49312b764ff33be3b9"
	var table RainbowTable
	f, err := os.Open(BIN_PATH)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&table); err != nil {
		panic(err)
	}
	password, found := Lookup(hashcode, table)
	fmt.Println(found)
	if found {
		fmt.Println("Password:", password)
	} else {
		fmt.Println("Password not found!!!")
	}
}
