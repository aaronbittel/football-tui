package term_utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	ClearLineCode        = "\033[0K"
	ClearScreenCode      = "\033[2J"
	HideCursorCode       = "\033[?25l"
	ShowCursorCode       = "\033[?25h"
	SaveCursorPosCode    = "\033[s"
	RestoreCursorPosCode = "\033[u"
	MoveCursorDownCode   = "\033[B"
	MoveCursorUpCode     = "\033[A"
	ResetCode            = "\033[m"
)

func ClearLineInst(pos ...int) string {
	var b strings.Builder
	switch len(pos) {
	case 1:
		b.WriteString(MoveCur(pos[0], 1))
	case 2:
		b.WriteString(MoveCur(pos[0], pos[1]))
	}
	b.WriteString(ClearLineCode)
	return b.String()
}

func MoveCur(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

func MoveCurLeft(count ...int) string {
	var b strings.Builder

	n := 1
	if len(count) > 0 {
		n = count[0]
	}
	b.WriteString(fmt.Sprintf("\033[%dD", n))

	return b.String()
}

func MoveCurDown() string {
	return MoveCurLeft() + fmt.Sprint(MoveCursorDownCode)
}

func MoveCurUp() string {
	return MoveCurLeft() + fmt.Sprint(MoveCursorUpCode)
}

var (
	ansiRegex = regexp.MustCompile("[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))")

	debugFunc = Debug()
)

func MoveCursorDown() {
	MoveCursorLeft()
	fmt.Print(MoveCursorDownCode)
}

func MoveCursorUp() {
	MoveCursorLeft()
	fmt.Print(MoveCursorUpCode)
}

func MoveCursorLeft(count ...int) {
	n := 1
	if len(count) > 0 {
		n = count[0]
	}
	fmt.Printf("\033[%dD", n)
}

func MoveCursorRight(count ...int) {
	n := 1
	if len(count) > 0 {
		n = count[0]
	}
	fmt.Printf("\033[%dC", n)
}

func HideCursor() {
	fmt.Print(HideCursorCode)
}

func ShowCursor() {
	fmt.Print(ShowCursorCode)
}

func ClearScreen() {
	fmt.Print(ClearScreenCode)
}

func ClearLine(pos ...int) {
	switch len(pos) {
	case 1:
		MoveCursor(pos[0], 1)
	case 2:
		MoveCursor(pos[0], pos[1])
	}
	fmt.Print(ClearLineCode)
}

func MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func GetSize(fd int) (rows, cols int, err error) {
	width, height, err := term.GetSize(fd)
	return height, width, err
}

func SaveCursorPos() {
	fmt.Print(SaveCursorPosCode)
}

func RestoreCursorPos() {
	fmt.Print(RestoreCursorPosCode)
}

func Debug() func(v ...any) {
	counter := 0
	times := 0
	debugRow := 39
	now := time.Now()
	return func(v ...any) {
		SaveCursorPos()
		defer RestoreCursorPos()

		if time.Since(now) > time.Second*3 {
			ClearLine(debugRow, 1)
			ClearLine(debugRow+1, 1)
			ClearLine(debugRow+2, 1)
			ClearLine(debugRow+3, 1)
			ClearLine(debugRow+4, 1)
			counter = 0
			now = time.Now()
		}

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
	}
}

func GetDebugFunc() func(v ...any) {
	return debugFunc
}

func TearDown() {
	fmt.Print(ResetCode)
	MoveCursor(1, 1)
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

func StringLen(str string) int {
	return utf8.RuneCountInString(StripAnsi(str))
}

// REFERENCE: https://github.com/acarl005/stripansi/blob/master/stripansi.go
func StripAnsi(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}

func WaitForEnter() {
	r := bufio.NewReader(os.Stdin)
	r.ReadByte()
}
