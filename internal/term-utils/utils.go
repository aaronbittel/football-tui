package term_utils

import (
	"fmt"
	"math/rand/v2"
	"syscall"

	"golang.org/x/term"
)

var debugFunc = Debug()

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

func ClearLine(pos ...int) {
	switch len(pos) {
	case 1:
		MoveCursor(pos[0], 1)
		fmt.Print("\033[0K")
	case 2:

		MoveCursor(pos[0], pos[1])
		fmt.Print("\033[0K")
	default:
		fmt.Print("\033[0K")
	}
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

func Debug() func(v ...any) {
	counter := 0
	times := 0
	debugRow := 40
	return func(v ...any) {
		SaveCursorPos()

		times = counter % 5

		if counter >= 5 {
			ClearLine(debugRow+times, 1)
		}

		out := ""
		for _, val := range v {
			out += fmt.Sprintf("%v ", val)
		}

		MoveCursor(debugRow+times, 1)
		fmt.Print(out)
		counter++
		RestoreCursorPos()
	}
}

func GetDebugFunc() func(v ...any) {
	return debugFunc
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

func getRandomColor() string {
	const colorCode = "\033[%dm"
	r := rand.IntN(8) + 30
	return fmt.Sprintf(colorCode, r)
}
