package term_utils

import "fmt"

const (
	Reset = "\033[0m"

	Orange    = "\033[38;5;208m"
	Blue      = "\033[38;5;12m"
	Green     = "\033[38;5;10m"
	Lightgray = "\033[38;5;240m"
	White     = "\033[37m"

	CtrlC = '\x03'
	Enter = 13
	Space = 32
	Tab   = 9

	FullBlock            = "█"
	WhiteSquare          = "▣"
	SquareDownHorizontal = "┬"
	SquareRightVertial   = "├"
	SquareLeftVertial    = "┤"
	SquareUpHorizontal   = "┴"
	SquareCross          = "┼"

	SquareTopLeft     = "┌"
	SquareTopRight    = "┐"
	SquareBottomLeft  = "└"
	SquareBottomRight = "┘"

	RoundedTopLeft     = "╭"
	RoundedTopRight    = "╮"
	RoundedBottomLeft  = "╰"
	RoundedBottomRight = "╯"

	HorizontalLine = "─"
	VerticalLine   = "│"
)

func Colorize(s, color string) string {
	return fmt.Sprintf("%s%s%s", color, s, Reset)
}
