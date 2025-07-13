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
				// capture and treat hashcode
				/*
					1. read hashcode from stdin
					2. decode hashcode
					3. return decoded hashcode
				*/
				// mustang codificado para sha256:
				hashcode := "a92f6bdb75789bccc118adfcf704029aa58063c604bab4fcdd9cd126ef9b69af"
				var rt []map[string]string
				arquivo, err := os.Open("rainbow_tables/jwt.secrets.bin")
				if err != nil {
					panic(err)
				}
				defer arquivo.Close()
				decoder := gob.NewDecoder(arquivo)
				if err := decoder.Decode(&rt); err != nil {
					panic(err)
				}
				_, ok := rt[0][hashcode]
				// nossa rainbow table ta quebrada, pau no nosso cu
				fmt.Println(ok)

				fmt.Println("opcao1")
			case 2:
				// read file and decode hashcodes trough the rainbow table
				// easy way: all in one thread:
				/* 1. open the file
				   2. fill a buffer with a few hashcode from the file
				   3. decode hashcodes from the buffer
				   4. print them out
				   5. if still are hashcodes remaining at the file -> back to step 2
				*/

				// cool parallel way:
				/*
					1. open the file for reading in goroutine A
					2. read hashcodes line by line (assuming one hashcode per line)
					3. publish each captured (hashcode, id) to a channel
						(the auto-incremented id helps preserve the original file order)
					4. for each (hashcode, id) received from the channel, launch a decode hashcode goroutine B
					5. each goroutine publishes its result as (decoded_hashcode, id) to another channel
					6. for each (decoded_hashcode, id) from the second channel, insert it into a binary tree ordered by id
					7. when all goroutines have finished, print the binary tree in order
				*/
				fmt.Println("opcao2")
			case 3:
				fmt.Println("CLOSING GORT, THANKS FOR USING!")
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
