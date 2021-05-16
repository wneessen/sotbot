package main

import (
	"fmt"
	"os"
)

func printHelp() {
	const usageText = `Usage:
    bot

Options:
    -h                   Show this help text`

	_, _ = fmt.Fprintf(os.Stderr, "%s\n", usageText)
	os.Exit(1)
}
