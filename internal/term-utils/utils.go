package term_utils

import "fmt"

func MoveCursorDown() {
	MoveCursorLeft()
	fmt.Print("\033[B")
}

func MoveCursorUp() {
	MoveCursorLeft()
	fmt.Print("\033[A")
}

func MoveCursorRight() {
	fmt.Print("\033[C")
}

func MoveCursorLeft() {
	fmt.Print("\033[D")
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func ClearScreen() {
	fmt.Print("\033[2J")
}

func MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func TearDown() {
	fmt.Print(Reset)
	MoveCursor(0, 0)
	ClearScreen()
	ShowCursor()
}

func Start() {
	HideCursor()
	ClearScreen()
	MoveCursor(1, 1)
}
