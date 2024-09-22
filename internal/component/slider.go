package component

import (
	"strings"
	term_utils "tui/internal/term-utils"
)

var cursors = []string{"┿", "▰"}

type Slider struct {
	cursor   string
	idx      int
	pos      int
	length   int
	row      int
	col      int
	vertical bool
}

func NewSlider(length int) *Slider {
	return &Slider{
		cursor: cursors[0],
		length: length,
		pos:    length / 2,
	}
}

func (s Slider) String() string {
	var b strings.Builder

	if s.vertical {
		b.WriteString("▄\n")
		for i := range s.length {
			if i == s.pos {
				b.WriteString(s.cursor + "\n")
			} else {
				b.WriteString(term_utils.VerticalLine + "\n")
			}
		}
		b.WriteString("▀\n")
	} else {
		b.WriteString("┠")
		for i := range s.length {
			if i == s.pos {
				b.WriteString(term_utils.FullBlock)
			} else {
				b.WriteString(term_utils.HorizontalLine)
			}
		}
		b.WriteString("┨")
	}

	return b.String()
}

func (s *Slider) Vertical() *Slider {
	s.vertical = false
	return s
}

func (s *Slider) MoveHigher() {
	if s.pos+1 >= s.length {
		return
	}
	s.pos++
}

func (s *Slider) MoveLower() {
	if s.pos-1 < 0 {
		return
	}
	s.pos--
	s.cursor = "┯"
}

func (s *Slider) NextCursor() {
	if s.idx+1 >= len(cursors) {
		s.idx = 0
	} else {
		s.idx++
	}
	s.cursor = cursors[s.idx]
}

func (s *Slider) At(row, col int) *Slider {
	s.row = row
	s.col = col
	return s
}

func (s Slider) Pos() (row, col int) {
	return s.row, s.col
}

func (s Slider) Lines() []string {
	return strings.Split(s.String(), "\n")
}
