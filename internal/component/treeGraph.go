package component

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"tui/internal/algorithm"
	term_utils "tui/internal/term-utils"
)

type Tree struct {
	instrCh chan<- string

	row int
	col int

	algo   *Algorithm
	frames []algorithm.ColumnGraphData
	Cursor int

	nums  []int
	level int
}

func NewTree(instrCh chan<- string, nums []int) *Tree {
	level := 0
	if len(nums) > 1 {
		level = int(math.Floor(math.Log2(float64(len(nums)))))
	}

	frames := make([]algorithm.ColumnGraphData, 0, 100)
	frames = append(frames, algorithm.NewColumnGraphData(slices.Clone(nums), map[int]string{}, ""))

	return &Tree{
		instrCh: instrCh,
		nums:    nums,
		level:   level,
		frames:  frames,
		row:     1,
		col:     1,
	}
}

func (t *Tree) AddFrame(next algorithm.ColumnGraphData) {
	prev := t.frames[t.Cursor]
	t.frames = append(t.frames, next)
	t.Cursor++
	t.PartialUpdate(prev, next)
}

func (t *Tree) Init(algo Algorithm) {
	columnCh := make(chan algorithm.ColumnGraphData)
	t.frames = algo.GetFrames(columnCh, slices.Clone(t.nums))
	t.Cursor = 0
	t.frames = slices.Clip(t.frames)
}

//            ╭───╮               ╭───╮         ╭────╮
//     	      │ 5 │               │ 5 │         │ 15 │
//            ╰─┬─╯               ╰─┬─╯         ╰─┬──╯
//       ╭──────┴──────╮        ╭───┴───╮     ╭───┴───╮
//     ╭─┴─╮         ╭─┴─╮    ╭─┴─╮   ╭─┴─╮ ╭─┴──╮ ╭──┴─╮
//     │ 7 │         │ 9 │    │ 7 │   │ 9 │ │ 17 │ │ 19 │
//     ╰─┬─╯         ╰─┬─╯    ╰───╯   ╰───╯ ╰────╯ ╰────╯
//   ╭───┴───╮     ╭───┴───╮
// ╭─┴──╮ ╭──┴─╮ ╭─┴──╮ ╭──┴─╮
// │ 13 │ │ 16 │ │ 10 │ │ 11 │
// ╰────╯ ╰────╯ ╰────╯ ╰────╯

//                          ╭───╮
//                          │ 5 │
//                          ╰─┬─╯
//              ╭─────────────┴─────────────╮
//            ╭─┴─╮                       ╭─┴─╮
//     	      │ 5 │                       │ 5 │
//            ╰─┬─╯                       ╰─┬─╯
//       ╭──────┴──────╮             ╭──────┴──────╮
//     ╭─┴─╮         ╭─┴─╮         ╭─┴─╮        ╭──┴─╮
//     │ 7 │         │ 9 │         │ 5 │        │ 15 │
//     ╰─┬─╯         ╰─┬─╯         ╰─┬─╯        ╰──┬─╯
//   ╭───┴───╮     ╭───┴───╮     ╭───┴───╮     ╭───┴───╮
// ╭─┴──╮ ╭──┴─╮ ╭─┴─╮  ╭──┴─╮ ╭─┴─╮   ╭─┴─╮ ╭─┴──╮  ╭─┴─╮
// │ 13 │ │ 16 │ │ 0 │  │ 11 │ │ 7 │   │ 9 │ │ 17 │  │ 9 │
// ╰────╯ ╰────╯ ╰───╯  ╰────╯ ╰───╯   ╰───╯ ╰────╯  ╰───╯

func (t Tree) PartialUpdate(prev, next algorithm.ColumnGraphData) string {
	var b strings.Builder
	pos := t.topPositions()
	for i := 0; i < len(prev.Nums()); i++ {
		// if nothing changed check the next number
		if prev.Nums()[i] == next.Nums()[i] && prev.Colors()[i] == next.Colors()[i] {
			continue
		}

		color := next.Colors()[i]

		nodeRelativePos := pos[i]
		nodeTopRow := t.row + nodeRelativePos[0]
		nodeTopCol := t.col + nodeRelativePos[1]

		b.WriteString(t.printNode(nodeTopRow, nodeTopCol, next.Nums(), i, color))
	}

	height := t.row + t.Height() + 1

	b.WriteString(term_utils.ClearLineInst(height, t.col))
	b.WriteString(term_utils.MoveCur(height, t.col))
	b.WriteString(next.Desc())

	return b.String()
}

func (t *Tree) Reset() {
	//TODO: HERE
	t.instrCh <- ""
}

func (t Tree) Pos() (row, col int) {
	return t.row, t.col
}

func (t *Tree) Next() {
	if t.Cursor+1 >= len(t.frames) {
		return
	}

	updateInstructions := t.PartialUpdate(t.frames[t.Cursor], t.frames[t.Cursor+1])
	t.Cursor++
	t.instrCh <- updateInstructions
}

func (t *Tree) Prev() {
	if t.Cursor-1 < 0 {
		return
	}

	updateInst := t.PartialUpdate(t.frames[t.Cursor], t.frames[t.Cursor-1])
	t.Cursor--
	t.instrCh <- updateInst
}

func (t *Tree) At(row, col int) {
	if row <= 0 {
		row = 1
	}
	if col <= 0 {
		col = 1
	}

	t.row = row
	t.col = col
}

func (t Tree) Chan() chan<- string {
	return t.instrCh
}

func (t Tree) String() string {
	var b strings.Builder

	hasSibling := func(i int) bool {
		return i+1 < len(t.nums)
	}

	positions := t.topPositions()

	rootNodePos := positions[0]
	b.WriteString(t.printNode(t.row+rootNodePos[0], t.col+rootNodePos[1], t.nums, 0))

	for i := 1; i < len(t.nums); i++ {
		relPos := positions[i]
		b.WriteString(t.printNode(t.row+relPos[0], t.col+relPos[1], t.nums, i))
	}

	connectorPositions := t.startingConnectorPositions()
	conn := 0
	for i := 1; i < len(t.nums); i += 2 {
		pos := connectorPositions[conn]
		length := pos[2]
		b.WriteString(term_utils.MoveCur(t.row+pos[0], t.col+pos[1]))
		b.WriteString(term_utils.RoundedTopLeft)
		b.WriteString(strings.Repeat(term_utils.HorizontalLine, length))

		if hasSibling(i) {
			b.WriteString(term_utils.SquareUpHorizontal)
			b.WriteString(strings.Repeat(term_utils.HorizontalLine, length))
			b.WriteString(term_utils.RoundedTopRight)
		} else {
			b.WriteString(term_utils.RoundedBottomRight)
		}
		conn++
	}

	b.WriteString(term_utils.ResetCode)
	return b.String()
}

func (t Tree) printNode(row, col int, nums []int, idx int, colors ...string) string {
	isLeft := func(i int) bool {
		return i != 0 && i%2 == 1
	}

	hasLeft := func(i int) bool {
		return i*2+1 < len(t.nums)
	}

	value := nums[idx]

	var b strings.Builder
	var helper strings.Builder

	helper.WriteString(term_utils.RoundedTopLeft + term_utils.HorizontalLine)

	if value >= 10 && !isLeft(idx) {
		helper.WriteString(term_utils.HorizontalLine)
	}

	if idx == 0 {
		helper.WriteString(term_utils.HorizontalLine)
	} else {
		helper.WriteString(term_utils.SquareUpHorizontal)
	}

	if value >= 10 && isLeft(idx) {
		helper.WriteString(term_utils.HorizontalLine)
	}

	helper.WriteString(term_utils.HorizontalLine + term_utils.RoundedTopRight + "\n")

	helper.WriteString(fmt.Sprintf("%s %d %s\n", term_utils.VerticalLine, value, term_utils.VerticalLine))

	helper.WriteString(term_utils.RoundedBottomLeft + term_utils.HorizontalLine)

	if value >= 10 && !isLeft(idx) {
		helper.WriteString(term_utils.HorizontalLine)
	}

	if hasLeft(idx) {
		helper.WriteString(term_utils.SquareDownHorizontal)
	} else {
		helper.WriteString(term_utils.HorizontalLine)
	}

	if value >= 10 && isLeft(idx) {
		helper.WriteString(term_utils.HorizontalLine)
	}

	helper.WriteString(term_utils.HorizontalLine + term_utils.RoundedBottomRight + "\n")

	if len(colors) >= 1 {
		b.WriteString(colors[0])
	}

	for i, line := range strings.Split(helper.String(), "\n") {
		if value < 10 {
			b.WriteString(term_utils.MoveCur(row+i, col-2))
		} else {
			if isLeft(idx) {
				b.WriteString(term_utils.MoveCur(row+i, col-2))
			} else {
				b.WriteString(term_utils.MoveCur(row+i, col-3))
			}
		}
		b.WriteString(line)
	}
	b.WriteString(term_utils.ResetCode)

	return b.String()
}

func (t Tree) Height() int {
	return (t.level+1)*3 + (t.level+1-1)*1
}

func (t Tree) Size() (rows, cols int) {
	//HACK: cols = ??

	// rows+2 to clear description
	return t.Height() + 2, 64
}

func (t Tree) topPositions() (topPositions [][]int) {
	switch t.level {
	case 0:
		return [][]int{{0, 2}}
	case 1:
		return [][]int{{0, 6}, {4, 2}, {4, 10}}
	case 2:
		return [][]int{{0, 13}, {4, 6}, {4, 20}, {8, 2}, {8, 10}, {8, 16}, {8, 24}}
	case 3:
		return [][]int{
			{0, 27},
			{4, 13}, {4, 41},
			{8, 6}, {8, 20}, {8, 34}, {8, 48},
			{12, 2}, {12, 10}, {12, 16}, {12, 24}, {12, 30}, {12, 38}, {12, 44}, {12, 52},
		}
	default:
		panic("not allowed tree with 4+ levels")
	}
}

func (t Tree) startingConnectorPositions() (connectorPositions [][]int) {
	switch t.level {
	case 0:
		return [][]int{}
	case 1:
		return [][]int{{3, 2, 3}}
	case 2:
		return [][]int{{3, 6, 6}, {7, 2, 3}, {7, 16, 3}}
	case 3:
		return [][]int{
			{3, 13, 13},
			{7, 6, 6}, {7, 34, 6},
			{11, 2, 3}, {11, 16, 3}, {11, 30, 3}, {11, 44, 3},
		}
	default:
		panic("not allowed tree with 4+ levels")
	}
}
