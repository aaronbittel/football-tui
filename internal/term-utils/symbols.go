package term_utils

import "fmt"

const (
	Reset = "\033[0m"

	Orange    = "\033[38;5;208m"
	Blue      = "\033[38;5;12m"
	Green     = "\033[38;5;10m"
	Lightgray = "\033[38;5;240m"
	White     = "\033[37m"

	BgRed        = "\033[48;2;255;95;95m"
	BgRedFgWhite = "\033[38;2;255;255;255;48;2;255;95;95m"

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

func Colorize(v any, color string) string {
	switch val := v.(type) {
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%s%d%s", color, val, Reset)
	case float32, float64:
		return fmt.Sprintf("%s%f%s", color, val, Reset)
	case string:
		return fmt.Sprintf("%s%s%s", color, val, Reset)
	default:
		return fmt.Sprintf("%s%v%s", color, val, Reset)
	}
}
