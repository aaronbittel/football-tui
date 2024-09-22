package component

import (
	"fmt"
	"slices"
	"strings"
	term_utils "tui/internal/term-utils"
)

var debug = term_utils.GetDebugFunc()

type ColumnGraphFrames struct {
	frames []ColumnGraphData
	cursor int

	originalNums []int

	row    int
	col    int
	height int
	width  int

	algo *Algorithm

	colParams columnParams
}

func NewColumnGraphFrames(nums []int) *ColumnGraphFrames {
	m := slices.Max(nums)
	spaces := len(fmt.Sprint(m))
	return &ColumnGraphFrames{
		frames:       make([]ColumnGraphData, 0, 100),
		originalNums: nums,
		height:       m + 1,
		width:        len(nums) + (len(nums)-1)*spaces,
		colParams: columnParams{
			maxVal: m,
			spaces: spaces,
		},
	}
}

func (cgf *ColumnGraphFrames) Idle() string {
	var b strings.Builder

	var char string

	for row := cgf.colParams.maxVal; row >= 1; row-- {
		for _, n := range cgf.originalNums {
			if n >= row {
				char = term_utils.FullBlock
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", cgf.colParams.spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range cgf.originalNums {
		s := cgf.colParams.spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	return b.String()
}

func (cgf ColumnGraphFrames) Nums() []int {
	return slices.Clone(cgf.originalNums)
}

func (cgf *ColumnGraphFrames) Init(algo *Algorithm) {
	columnCh := make(chan ColumnGraphData)

	go algo.AlgorithmFn(columnCh, slices.Clone(cgf.originalNums))

	// Resetting from previous visualization
	cgf.frames = make([]ColumnGraphData, 0, 100)
	cgf.cursor = 0
	for col := range columnCh {
		cgf.frames = append(cgf.frames, col)
	}

	cgf.algo = algo
}

func (cgf ColumnGraphFrames) partialUpdate(prev, next ColumnGraphData) {
	for i := 0; i < len(prev.nums); i++ {
		// Only update the columns that have changed
		if prev.nums[i] == next.nums[i] && prev.colors[i] == next.colors[i] {
			continue
		}

		cgf.removeColumn(prev.nums[i], i)
		cgf.printNewCol(next.nums[i], i, next.colors[i])
		cgf.updateDescription(prev.desc, next.desc)

	}
}

func (cgf *ColumnGraphFrames) Next() {
	if cgf.cursor+1 >= len(cgf.frames) {
		return
	}

	cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor+1])
	cgf.cursor++
}

func (cgf *ColumnGraphFrames) Prev() {
	if cgf.cursor-1 < 0 {
		return
	}

	cgf.partialUpdate(cgf.frames[cgf.cursor], cgf.frames[cgf.cursor-1])
	cgf.cursor--
}

func (cgf ColumnGraphFrames) Print() {
	var b strings.Builder

	var char string
	frame := cgf.frames[cgf.cursor]

	for row := cgf.colParams.maxVal; row >= 1; row-- {
		for i, n := range frame.nums {
			if n >= row {
				if color, found := frame.colors[i]; found {
					char = fmt.Sprint(color, term_utils.FullBlock, term_utils.Reset)
				} else {
					char = term_utils.FullBlock
				}
			} else {
				char = " "
			}
			b.WriteString(char + strings.Repeat(" ", cgf.colParams.spaces))
		}
		b.WriteString("\n")
	}

	for _, n := range frame.nums {
		s := cgf.colParams.spaces - len(fmt.Sprint(n)) + 1
		b.WriteString(fmt.Sprintf("%d%s", n, strings.Repeat(" ", s)))
	}

	width := len(cgf.originalNums) + (len(cgf.originalNums)-1)*cgf.colParams.spaces

	b.WriteString("\n")
	leftPad := (width - StringLen(frame.desc)) / 2
	if leftPad > 0 {
		b.WriteString(strings.Repeat(" ", leftPad))
	}
	b.WriteString(frame.desc)

	for i, line := range strings.Split(b.String(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(line)
	}
}

func (cgf *ColumnGraphFrames) At(row, col int) *ColumnGraphFrames {
	cgf.row = row
	cgf.col = col
	return cgf
}

func (cgf ColumnGraphFrames) PrintIdle() {
	for i, line := range strings.Split(cgf.Idle(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(line)
	}
}

func (cgf ColumnGraphFrames) Clear() {
	width := len(cgf.originalNums) + (len(cgf.originalNums)-1)*cgf.colParams.spaces
	for i := range strings.Split(cgf.Idle(), "\n") {
		term_utils.MoveCursor(cgf.row+i, cgf.col)
		fmt.Print(strings.Repeat(" ", width))
	}
}

func (cgf ColumnGraphFrames) moveCursorTopColumn(colHeight, colIdx int) {
	term_utils.MoveCursor(cgf.row+(cgf.height-colHeight-1), cgf.col+colIdx*(cgf.colParams.spaces+1))
}

func (cgf ColumnGraphFrames) removeColumn(colHeight, colIdx int) {
	// -1 because the number underneath the columns counts into the cgf.height
	cgf.moveCursorTopColumn(colHeight, colIdx)
	for range colHeight {
		fmt.Print(" ")
		term_utils.MoveCursorDown()
	}
	fmt.Print("  ")

	term_utils.MoveCursorLeft()
	term_utils.MoveCursorLeft()
}

func (cgf ColumnGraphFrames) printNewCol(colHeight, colIdx int, color string) {
	fmt.Print(colHeight)

	// Move cursor to first column segment
	term_utils.MoveCursor(cgf.row+cgf.height-2, cgf.col+colIdx*(cgf.colParams.spaces+1))
	fmt.Print(color)
	for range colHeight {
		fmt.Print(term_utils.FullBlock)
		term_utils.MoveCursorUp()
	}
	fmt.Print(term_utils.Reset)
}

func (cgf ColumnGraphFrames) updateDescription(prev, next string) {
	leftPad := (cgf.width - StringLen(next)) / 2
	term_utils.MoveCursor(cgf.DescCol(), cgf.col)

	fmt.Print(strings.Repeat(" ", leftPad))
	fmt.Print(next)

	diff := StringLen(prev) - StringLen(next)
	if diff > 0 {
		fmt.Print(strings.Repeat(" ", (diff/2)+1))
	}
}

func (cgf ColumnGraphFrames) ClearDescription() {
	frame := cgf.frames[cgf.cursor]

	leftPad := (cgf.width - StringLen(frame.desc)) / 2
	debug("POS", cgf.DescCol(), cgf.col+leftPad)
	term_utils.MoveCursor(cgf.DescCol(), cgf.col+leftPad)
	fmt.Print(strings.Repeat(" ", StringLen(frame.desc)))
}

func (cgf ColumnGraphFrames) DescCol() int {
	return cgf.row + cgf.height + 1
}
