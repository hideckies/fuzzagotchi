package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/fuzzer"
	"github.com/hideckies/fuzzagotchi/pkg/output"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		select {
		case <-sigCh:
			fmt.Println("Keyboard interrupt detected, terminating.")
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

	// Count the number of words
	wordlist, err := os.ReadFile(cmd.Options.Wordlist)
	if err != nil {
		panic(err)
	}

	totalWords := len(strings.Split(string(wordlist), "\n"))

	// Create a new Fuzzer and start fuzzing
	fuzzer := fuzzer.NewFuzzer(ctx, cmd.Options, fuzztype, totalWords)
	fuzzer.Run()
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
