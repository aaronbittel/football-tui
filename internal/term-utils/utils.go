package term_utils

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

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

func GetSize(fd int) (rows, cols int, err error) {
	return term.GetSize(fd)
}

func SaveCursorPos() {
	fmt.Print("\033[s")
}

func RestoreCursorPos() {
	fmt.Print("\033[u")
}

func Debug(v ...any) {
	SaveCursorPos()
	MoveCursor(40, 1)
	out := ""
	for _, val := range v {
		out += fmt.Sprintf("%v ", val)
	}
	fmt.Print(out)
	RestoreCursorPos()
}

func TearDown() {
	fmt.Print(Reset)
	MoveCursor(0, 0)
	ClearScreen()
	ShowCursor()
}

func Start() (fd int, oldState *term.State, err error) {
	HideCursor()
	ClearScreen()
	MoveCursor(1, 1)

	fd = int(syscall.Stdin)
	oldState, err = term.MakeRaw(fd)

	return fd, oldState, err
}
