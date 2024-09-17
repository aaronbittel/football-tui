package term_utils

import "fmt"

const (
	Reset = "\033[0m"

	Orange = "\033[38;5;208m"
	Blue   = "\033[38;5;12m"
	Green  = "\033[38;5;10m"

	CtrlC   = '\x03'
	ArrowUp = '\n'

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
