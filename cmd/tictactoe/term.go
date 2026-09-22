package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func setupTerm() func() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, "This program must be run from an interactive terminal")
		os.Exit(1)
	}

	termState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not enter raw mode: %e", err)
	}

	enterAltScreen()
	cursorHidden()

	return func() {
		cursorVisible()
		exitAltScreen()
		term.Restore(int(os.Stdin.Fd()), termState)
	}
}
