package main

import (
	"bufio"
	"crypto/sha512"
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

	for opt != 5 {
		// cleaning stdout buffer with ansi scape code
		fmt.Print("\033[H\033[2J")
		switch opt {
		// capture the entries with []byte always, there are lighter than string.
		// and the https://pkg.go.dev/crypto (native go lib for crypto hashing) works by default with []byte;
		case 1:
			GeneratePasswords(PASSWORDS_PATH, PASSWORD_LENGTH, NUM_PASSWORDS)
		case 2:
			GenerateTable("", BIN_PATH)
		case 3:
			DecryptHashcode()
		case 4:
			TestRandomPasswords(5)
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
	time.Sleep(time.Millisecond * 250)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 250)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 250)
	fmt.Print(".")
	time.Sleep(time.Millisecond * 250)
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
	fmt.Print("Enter the hashcode to test: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	hashcode := scanner.Text()
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

func TestRandomPasswords(passwordNumber int) {
	fmt.Println("Testing", passwordNumber, "random passwords against the rainbow table...")
	var table RainbowTable
	f, err := os.Open(BIN_PATH)
	if err != nil {
		fmt.Println("Could not open rainbow table file:", err)
		return
	}
	defer f.Close()
	decoder := gob.NewDecoder(f)
	if err := decoder.Decode(&table); err != nil {
		fmt.Println("Could not decode rainbow table:", err)
		return
	}
	hasher := sha512.New()
	for i := 1; i <= passwordNumber; i++ {
		pw := GeneratePassword(PASSWORD_LENGTH)
		hasher.Reset()
		hasher.Write([]byte(pw))
		hash := hasher.Sum(nil)
		hashHex := fmt.Sprintf("%x", hash)
		fmt.Printf("%d) Password: %s\n   Hash: %s\n", i, pw, hashHex)
		foundPw, found := Lookup(hashHex, table)
		if found {
			fmt.Printf("   FOUND in rainbow table: %s\n", foundPw)
		} else {
			fmt.Println("   NOT FOUND in rainbow table.")
		}
	}
}
