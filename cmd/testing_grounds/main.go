package main

import (
	"time"
	term_utils "tui/internal/term-utils"

	"golang.org/x/term"
)

func main() {
	time.Sleep(time.Second)
	fd, oldState, err := term_utils.Start()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	defer term_utils.TearDown()

	fixingRaceCondition()
}
