package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"

	"golang.org/x/term"
)

const (
	reset  = "\033[0m"
	orange = "\033[38;5;208m"
	ctrlC  = '\x03'
)

type pos struct {
	row int
	col int
}

func NewPos(row, col int) pos {
	return pos{
		row: row,
		col: col,
	}
}

func main() {
	start()
	fd := int(syscall.Stdin)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)

	// contentPos := NewPos(4, 20)

	reader := bufio.NewReader(os.Stdin)

	createContent := func(i, j int) {
		clear()
		tabsPos := NewPos(0, 20)
		listPos := NewPos(5, 0)
		tabs := NewTabs("t table", "f fixtures & results", "s stats", "p table predicter").SetTab(i)
		list := NewList("Bundesliga", "2. Bundesliga", "Premier League", "La Liga", "Ligue 1", "Serie A").Select(j)

		Print(tabsPos, tabs.String())
		Print(listPos, list.String())
	}

	listSelected := 0
	tabsSelected := 0
	contentPos := NewPos(4, 20)
	go table(contentPos)

	createContent(tabsSelected, listSelected)

outer:
	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		moveCursor(30, 50)
		fmt.Print(listSelected, tabsSelected, string(b))

		switch b {
		case 'q', ctrlC:
			break outer
		case 't':
			if tabsSelected == 0 {
				continue
			}
			tabsSelected = 0
			createContent(tabsSelected, listSelected)
			go table(contentPos)
		case 'f':
			if tabsSelected == 1 {
				continue
			}
			tabsSelected = 1
			createContent(tabsSelected, listSelected)
			go matchday(contentPos)
		case 's':
			if tabsSelected == 2 {
				continue
			}
			tabsSelected = 2
			createContent(tabsSelected, listSelected)
		case 'p':
			if tabsSelected == 3 {
				continue
			}
			tabsSelected = 3
			createContent(tabsSelected, listSelected)
		case 'k':
			if listSelected-1 >= 0 {
				listSelected--
				createContent(tabsSelected, listSelected)
			}
		case 'j':
			if listSelected+1 <= 5 {
				listSelected++
				createContent(tabsSelected, listSelected)
			}
		}
	}

	tearDown()
}

func firstVersion() {
	start()
	defer tearDown()
	fd := int(syscall.Stdin)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)

	boxMatchday := NewBox("Matchday").WithPadding().WithRoundedCorners().WithColoredBorder(orange)
	boxTable := NewBox("Table").WithPadding().WithRoundedCorners()
	Print(NewPos(5, 40), boxMatchday.String())
	Print(NewPos(5, 55), boxTable.String())
	go matchday(NewPos(10, 20))

	selected := 0

	reader := bufio.NewReader(os.Stdin)
	Print(NewPos(0, 0), "Press 'q' or 'ctrl-c' to quit.")

outer:
	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', ctrlC:
			break outer
		case 'm':
			if selected == 0 {
				break
			}
			selected = 0
			boxMatchday = NewBox("Matchday").WithPadding().WithRoundedCorners().WithColoredBorder(orange)
			boxTable = NewBox("Table").WithPadding().WithRoundedCorners()
			go matchday(NewPos(10, 20))
		case 't':
			if selected == 1 {
				break
			}
			selected = 1
			boxMatchday = NewBox("Matchday").WithPadding().WithRoundedCorners()
			boxTable = NewBox("Table").WithPadding().WithRoundedCorners().WithColoredBorder(orange)
			go table(NewPos(10, 20))
		}
		clear()
		Print(NewPos(5, 40), boxMatchday.String())
		Print(NewPos(5, 55), boxTable.String())
	}
}
func table(pos pos) {
	tableInfo := getTable()
	table := NewTable(
		NewHeader("#", true),
		NewHeader("TeamName", false),
		NewHeader("Games", true),
		NewHeader("W", true),
		NewHeader("D", true),
		NewHeader("L", true),
		NewHeader("Goals", true),
		NewHeader("Diff", true),
		NewHeader("Points", true),
	).WithRoundedCorners()

	for i, info := range tableInfo {
		data := fmt.Sprintf("%d??%s??%d??%d??%d??%d??%d??%d??%d", i+1, info.TeamName, info.Matches, info.Won, info.Draw, info.Lost, info.Goals, info.GoalDiff, info.Points)
		table.AddRow(strings.Split(data, "??"))
	}

	Print(pos, table.String())
}

func matchday(pos pos) {
	matches := getMatchday(1)
	matchStrs := make([][]string, 0, 9)
	maxLenHome := 0
	for _, match := range matches {
		var (
			result    = match.MatchResults[1]
			home      = match.Home.TeamName
			away      = match.Away.TeamName
			homeGoals = result.PointsHome
			awayGoals = result.PointsAway
		)
		if utf8.RuneCountInString(home) > maxLenHome {
			maxLenHome = utf8.RuneCountInString(home)
		}

		matchStrs = append(matchStrs, []string{home, strconv.Itoa(homeGoals), strconv.Itoa(awayGoals), away})
	}

	content := make([]string, 9, 9)
	for i, c := range matchStrs {
		var (
			home      = c[0]
			homeGoals = c[1]
			awayGoals = c[2]
			away      = c[3]
		)
		b := new(strings.Builder)

		b.WriteString(strings.Repeat(" ", maxLenHome-utf8.RuneCountInString(home)))
		b.WriteString(fmt.Sprintf("%s %s : %s %s", home, homeGoals, awayGoals, away))

		content[i] = b.String()
	}

	box := NewBox(content...).WithRoundedCorners().WithTitle("Matchday 2").WithPadding(2, 4, 1, 2)
	Print(pos, box.String())
}

func terminalTesting() {
	start()
	defer tearDown()
	fd := int(syscall.Stdin)

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Error setting raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)

	box := NewBox("").WithTitle("Title").WithPadding(4, 4, 4, 4).WithRoundedCorners()
	Print(NewPos(0, 0), box.String())

	reader := bufio.NewReader(os.Stdin)
	Print(NewPos(50, 10), "Press 'q' to quit.")

	boxes := createBoxes(0)
	printBoxes(10, 10, boxes)

	var i int
outer:
	for {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', ctrlC:
			break outer
		case '0':
			i = 0
		case '1':
			i = 1
		case '2':
			i = 2
		case '3':
			i = 3
		}

		boxes := createBoxes(i)
		printBoxes(10, 10, boxes)
	}
}

func printBoxes(row, col int, boxes []string) {
	vSpace := 7
	hSpace := 10
	Print(NewPos(row, col), boxes[0])
	Print(NewPos(row, col+hSpace), boxes[1])
	Print(NewPos(row+vSpace, col+hSpace), boxes[2])
	Print(NewPos(row+vSpace, col), boxes[3])
}

func Print(pos pos, str string) {
	if pos.row == 0 {
		pos.row = 1
	}
	if pos.col == 0 {
		pos.col = 1
	}
	for i, s := range strings.Split(str, "\n") {
		moveCursor(pos.row+i, pos.col)
		fmt.Print(s)
	}
}

func createBoxes(focused int) []string {
	boxes := []string{}
	for i := range 4 {
		box := NewBox(fmt.Sprintf("%d", i)).WithPadding().WithCenteredText().WithRoundedCorners()

		if i == focused {
			box.WithColoredBorder(getRandomColor())
		}

		boxes = append(boxes, box.String())
	}
	return boxes
}

func getRandomColor() string {
	const colorCode = "\033[%dm"
	r := rand.IntN(8) + 30
	return fmt.Sprintf(colorCode, r)
}

func start() {
	hideCursor()
	clear()
	moveCursor(1, 1)
}

func moveCursorDown() {
	fmt.Print("\033[B")
}

func moveCursorUp() {
	fmt.Print("\033[A")
}

func moveCursorRight() {
	fmt.Print("\033[C")
}

func moveCursorLeft() {
	fmt.Print("\033[D")
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func clear() {
	fmt.Print("\033[2J")
}

func moveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func tearDown() {
	fmt.Print(reset)
	moveCursor(0, 0)
	clear()
	showCursor()
}

func StringLen(str string) int {
	return utf8.RuneCountInString(StripAnsi(str))
}

func Mask(str string) []int {
	mask := []int{}
	for _, l := range strings.Split(str, "\n") {
		mask = append(mask, StringLen(l))
	}
	return mask
}

func waitForEnter() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
