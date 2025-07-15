package cmd

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"main/utils"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Rainbow table hashcode decoding application",
	Long: `
		Rainbow Table implementation:
			- Alphabet length = 70 symbols. For example;
			- It's a rainbow table full of colors (different hash functions) not a monochromatic one;
			- It supports passwords with 5 or more characters;
			- Focused on SHA512 hash;
		`,
	Run: func(cmd *cobra.Command, args []string) {
		// printing menu
		fmt.Print(ASCII)
		fmt.Print(MENU)

		// scanning first line
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		opt, _ := utils.BytesToOption(scanner.Bytes())

		for opt != 3 {
			// cleaning stdout buffer with ansi scape code
			fmt.Print("\033[H\033[2J")
			switch opt {
			// capture the entries with []byte always, there are lighter than string.
			// and the https://pkg.go.dev/crypto (native go lib for crypto hashing) works by default with []byte;
			case 1:
				hashcode := "bdb35b4df1ccd8e0f3d113eb6840472f1702f397820a716dc26b419485bd62e3d8eebd639dec11a998990f69bac58fdf4cebf1abea2f00fca62ca68f1d195e1d"
				var rt [][]utils.ChainEntry
				arquivo, err := os.Open("rainbow_tables/secrets2.bin")
				if err != nil {
					panic(err)
				}
				defer arquivo.Close()
				decoder := gob.NewDecoder(arquivo)
				if err := decoder.Decode(&rt); err != nil {
					panic(err)
				}
				password, found := utils.LookupInRainbowTable(hashcode, rt[1], utils.GetSHA512(), CHAIN_SIZE)
				fmt.Println(found)
				if found {
					fmt.Println("Password:", password)
				} else {
					fmt.Println("Password not found")
				}
				fmt.Println("opcao1")
			case 2:
				fmt.Println("opcao2")
			case 3:
				fmt.Println("CLOSING GORT, THANKS FOR USING!")
			case 4:
				fmt.Println("Digite a senha para testar:")
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				input := scanner.Text()
				fmt.Println("Senha inserida:", input)

				chainLen := CHAIN_SIZE
				h := utils.GetSHA512()
				plaintext := input
				var end string
				fmt.Println("\n--- Gerando chain ---")
				for i := 0; i < chainLen; i++ {
					h.Reset()
					h.Write([]byte(plaintext))
					hashcode := h.Sum(nil)
					hashHex := utils.EncodeHex(hashcode)
					fmt.Printf("[%d] Hash: %s\n", i, hashHex)
					reduced := utils.Reduce(hashcode, i)
					fmt.Printf("[%d] Reduced: %s\n", i, reduced)
					plaintext = reduced
					if i == chainLen-1 {
						end = reduced
					}
				}
				fmt.Println("--- Fim da chain ---")
				fmt.Printf("Entrada gerada na rainbow table: (%s|%s)\n\n", input, end)
			default:
				// continue asking for input
				fmt.Println("error: invalid option")
			}
			// this is only because we're cleaning the console with ansi scape chars
			time.Sleep(time.Second * 1)
			fmt.Print(ASCII)
			fmt.Print(MENU)
			scanner.Scan()
			opt, _ = utils.BytesToOption(scanner.Bytes())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
