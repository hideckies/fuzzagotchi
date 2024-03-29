package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/fuzzer"
	"github.com/hideckies/fuzzagotchi/pkg/output"
)

// Main function to start running Fuzzagotchi.
func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	var f fuzzer.Fuzzer
	var err error

	go func() {
		select {
		case <-sigCh:
			fmt.Println("\x1b[2KKeyboard interrupt detected, terminating.")
			f.Scan()
			cancel()
			os.Exit(0)
		case <-ctx.Done():
			return
		}
	}()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !cmd.Proceed {
		return
	}

	// The "color" package setting
	color.NoColor = cmd.Options.NoColor

	output.Banner(cmd.Options)

	// Detect the fuzz type
	fuzztype := detectFuzzType(cmd.Options)

	// Check if the -w flag is buildin list.
	wordlistType := ""
	reAlpha := regexp.MustCompile(`ALPHA_[A-Z]+_[A-Z]+`)
	reNum := regexp.MustCompile(`NUM_[0-9]+_[0-9]+`)
	totalWords := 0
	if reAlpha.MatchString(cmd.Options.Wordlist) || reNum.MatchString(cmd.Options.Wordlist) {
		wordlistType = "builtin"
	} else {
		// Count the number of words
		wordlist, err := os.ReadFile(cmd.Options.Wordlist)
		if err != nil {
			panic(err)
		}
		totalWords = len(strings.Split(string(wordlist), "\n"))
	}

	// Create a new Fuzzer and start fuzzing
	f, err = fuzzer.NewFuzzer(ctx, cmd.Options, fuzztype, wordlistType, totalWords)
	if err != nil {
		color.Red("%s", err)
		return
	}
	if err = f.Run(); err != nil {
		fmt.Printf("%v", err)
	}
}

// Detect the fuzz type
func detectFuzzType(opts cmd.CmdOptions) string {
	vals := reflect.ValueOf(opts)
	types := vals.Type()
	for i := 0; i < vals.NumField(); i++ {
		key := types.Field(i).Name
		val := vals.Field(i).String()
		if strings.Contains(val, "EGG") {
			return key
		}
	}
	return ""
}
