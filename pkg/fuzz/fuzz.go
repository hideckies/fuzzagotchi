package fuzz

import (
	"bufio"
	"os"

	"github.com/fatih/color"
)

func Fuzz(url string, verbose bool, wordlist string) {
	readFile, err := os.Open(wordlist)

	if err != nil {
		color.HiRed("%v\nPlease install seclists by running 'sudo apt install seclists'.\n", err)
		os.Exit(0)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	var word string
	for fileScanner.Scan() {
		word = fileScanner.Text()

		_ = word
	}

	readFile.Close()
}
