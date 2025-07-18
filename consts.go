package main

// RAINBOW TABLE CONSTS
const (
	CHARSET         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#%^&*?"
	CHAIN_SIZE      = 1000
	NUM_PASSWORDS   = 3870000
	PASSWORD_LENGTH = 5

	// MENU CONSTS
	ASCII = `
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
  ▄████  ▄█████  ▄██▀▀▀█▄ ██████▓
 ██▒ ▀█ ▒██▒  ██ ▓█     █   ▓██   
 ██  ▄▄ ▒██░  ██ ▓█   ▄█    ▓██  
 ▓█   █ ▒██   ██ ▒█▀▀█▄     ▓█▓
 ░▓███▀  ░████▓  ░█   ▒█   ▄█▓▓
▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
BY> Gustavo Ventino & Arthur Faria
# Rainbow Table implementation in Go.
`
	MENU = `
Choose an option from below:
	1. Generate Passwords
	2. Generate Table
	3. Try to decrypt hashcode
	4. Test 5 random passwords in rainbow table
	5. Close

input: `

	MAX_INPUT_SIZE  = 20
	MAX_OPTION_SIZE = 2

	PASSWORDS_BUFF_SIZE = 100000
	PASSWORDS_PATH      = "passwords/secrets.txt"
	BIN_PATH            = "tables/rainbow_table.bin"
)
