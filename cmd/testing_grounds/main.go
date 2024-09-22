package main

import (
	"time"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

func main() {
	time.Sleep(time.Second)
	fd, oldState, err := utils.Start()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	defer utils.TearDown()

	fixingRaceCondition()
}
