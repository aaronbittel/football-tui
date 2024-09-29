package main

import (
	term_utils "tui/internal/term-utils"

	"golang.org/x/term"
)

func main() {
	fd, oldState, err := term_utils.Start()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	defer term_utils.TearDown()
}
