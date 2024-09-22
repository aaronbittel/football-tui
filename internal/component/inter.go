package component

import (
	"fmt"
	"regexp"
	utils "tui/internal/term-utils"
	"unicode/utf8"
)

type Masker interface {
	Lines() []string
}

type Updater interface {
	Clearer
	Printer
}

type Clearer interface {
	Masker
	Pos() (row, col int)
}

type Printer interface {
	Masker
	Pos() (row, col int)
}

func StringLen(str string) int {
	return utf8.RuneCountInString(StripAnsi(str))
}

// REFERENCE: https://github.com/acarl005/stripansi/blob/master/stripansi.go
func StripAnsi(str string) string {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	re := regexp.MustCompile(ansi)
	return re.ReplaceAllString(str, "")
}

func Mask(m Masker) (height, width int) {
	lines := m.Lines()
	return len(lines), StringLen(lines[0])
}

func Print(s Printer) {
	row, col := s.Pos()
	for i, s := range s.Lines() {
		utils.MoveCursor(row+i, col)
		fmt.Print(s)
	}
}

func Update(u Updater) {
	Clear(u)
	Print(u)
}
