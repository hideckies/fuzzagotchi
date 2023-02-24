package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

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

	// Create a new Fuzzer and start fuzzing
	fuzzer := fuzzer.NewFuzzer(cmd.Options, ctx)
	fuzzer.Run()
}
