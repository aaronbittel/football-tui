package main

import (
	"bufio"
	"fmt"
	"os"
	component "tui/internal/component"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

func main() {
	fd, oldState, err := utils.Start()
	if err != nil {
		panic(err)
	}
	defer term.Restore(fd, oldState)
	defer utils.TearDown()

	slider := component.NewSlider(25).At(10, 10)
	component.Print(slider)

	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', utils.CtrlC:
			running = false
		case 'n':
			slider.NextCursor()
			component.Clear(slider)
			component.Print(slider)
		case 'u':
			slider.MoveHigher()
			component.Clear(slider)
			component.Print(slider)
		case 'd':
			slider.MoveLower()
			component.Clear(slider)
			component.Print(slider)
		}
	}
}
