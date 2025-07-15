package cmd

const (
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
	1. Try to decrypt one hashcode from stdin;
	2. Try to decrypt several hashcodes from file;
	3. Close;
	4. Test chain gen;
input: `

	MAX_INPUT_SIZE = 20

	// len(option) + len(\n) = 1 + 1 = 2
	MAX_OPTION_SIZE = 2

	CHAIN_SIZE          = 50
	PASSWORDS_BUFF_SIZE = 10000
)
