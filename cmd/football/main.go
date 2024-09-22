package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"tui/internal/component"
	utils "tui/internal/term-utils"

	"golang.org/x/term"
)

type Match struct {
	Date      time.Time
	Home      string
	Away      string
	HomeGoals int
	AwayGoals int
}

func NewMatch(date time.Time, home, away string, homeGoals, awayGoals int) Match {
	return Match{
		Date:      date,
		Home:      home,
		Away:      away,
		HomeGoals: homeGoals,
		AwayGoals: awayGoals,
	}
}

func (m Match) String() string {
	return fmt.Sprintf("%s %d : %d %s", m.Home, m.HomeGoals, m.AwayGoals, m.Away)
}

func main() {
	fd, oldState, err := utils.Start()
	if err != nil {
		log.Fatal("error initalizing terminal")
	}
	defer term.Restore(fd, oldState)
	defer utils.TearDown()

	matchday := 1
	matches := LoadFromCSV("bundesliga", 2023, matchday)

	matchBox := component.NewBox(
		createMatchdayStrings(matches)...).
		WithTitle(fmt.Sprintf("Matchday %d", matchday)).
		WithRoundedCorners().
		WithPadding(1, 5).
		WithColoredBorder(utils.Blue).
		At(10, 10)

	controlBox := component.NewBox(
		"[ n ] - next", "[ p ] - previous").
		WithRoundedCorners().
		WithTitle("Controls").
		WithPadding(0, 5).
		At(24, 25)

	component.Print(matchBox)
	component.Print(controlBox)

	curMatchday := 1

	reader := bufio.NewReader(os.Stdin)
	running := true

	for running {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println("Error reading byte:", err)
			break
		}

		switch b {
		case 'q', utils.CtrlC:
			running = false
		case 'n':
			if curMatchday+1 > 34 {
				break
			}
			curMatchday++
			matches := LoadFromCSV("bundesliga", 2023, curMatchday)
			matchBox.Update(fmt.Sprintf("Matchday %d", curMatchday), createMatchdayStrings(matches)...)
			utils.Debug(curMatchday)
		case 'p':
			if curMatchday-1 < 0 {
				break
			}
			curMatchday--
			matches := LoadFromCSV("bundesliga", 2023, curMatchday)
			matchBox.Update(fmt.Sprintf("Matchday %d", curMatchday), createMatchdayStrings(matches)...)
			utils.Debug(curMatchday)
		}
	}
}

func createMatchdayStrings(matches []Match) []string {
	matchStrs := make([]string, 9, 9)

	var longestHome int
	for _, m := range matches {
		if utf8.RuneCountInString(m.Home) > longestHome {
			longestHome = utf8.RuneCountInString(m.Home)
		}
	}

	for i, m := range matches {
		space := longestHome - utf8.RuneCountInString(m.Home)
		matchStrs[i] = strings.Repeat(" ", space) + m.String()
	}

	return matchStrs
}

func LoadFromCSV(league string, season, matchday int) []Match {
	var matches []Match
	str := func(i int) string {
		return fmt.Sprintf("%d", i)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error getting home dir", err)
	}
	path := path.Join(home, "projects", "golang", "scraper", "data", league, str(season))

	filename := fmt.Sprintf("%s/%d.csv", path, matchday)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("error opening file", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	_, err = reader.Read()
	if err != nil {
		log.Fatal("error reading csv header", err)
	}

	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatal("error reading csv header", err)
	}

	layout := "2006-01-02 15:04:05 -0700 MST"

	var m Match
	for _, row := range rows {
		date, err := time.Parse(layout, row[4])
		if err != nil {
			log.Fatal("error parsing date", err)
		}
		homeGoals, err := strconv.Atoi(row[1])
		if err != nil {
			log.Fatal("error converting number", err)
		}
		awayGoals, err := strconv.Atoi(row[2])
		if err != nil {
			log.Fatal("error converting number", err)
		}
		m = NewMatch(date, row[0], row[3], homeGoals, awayGoals)
		matches = append(matches, m)
	}

	return matches
}
