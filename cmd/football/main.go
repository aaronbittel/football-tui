package main

func main() {
	// firstVersion()
	// secondVersion()
}

// func secondVersion() {
// 	component.Start()
// 	fd := int(syscall.Stdin)
//
// 	oldState, err := term.MakeRaw(fd)
// 	if err != nil {
// 		fmt.Println("Error setting raw mode:", err)
// 		return
// 	}
// 	defer term.Restore(fd, oldState)
//
// 	reader := bufio.NewReader(os.Stdin)
//
// 	listSelected := 0
// 	tabsSelected := 1
// 	selectedMatchday := 2
// 	colored := false
//
// 	createContent := func(i, j int) {
// 		component.ClearScreen()
// 		tabs := component.NewTabs("t table", "f fixtures & results", "s stats", "p table predicter").Select(i).At(1, 20)
// 		list := component.NewList("Bundesliga", "2. Bundesliga", "3. Bundesliga", "DFB Pokal", "Champions League", "Europa League").Select(j).At(5, 1)
//
// 		component.Print(tabs)
// 		component.Print(list)
//
// 		if tabsSelected == 1 {
// 			box := component.NewBox(fmt.Sprintf("Matchday: %d", selectedMatchday)).
// 				WithRoundedCorners().
// 				WithTitle("Settings").
// 				WithPadding(1, 1, 4, 1).
// 				At(4, 83)
// 			if colored {
// 				box.WithColoredBorder(component.Orange)
// 			}
// 			component.Print(box)
// 		}
// 	}
//
// 	go matchday(4, 20, 2)
//
// 	createContent(tabsSelected, listSelected)
//
// outer:
//
// 	for {
// 		b, err := reader.ReadByte()
// 		if err != nil {
// 			fmt.Println("Error reading byte:", err)
// 			break
// 		}
//
// 		switch b {
// 		case 'q', ctrlC:
// 			break outer
// 		case 't':
// 			if tabsSelected == 0 {
// 				continue
// 			}
// 			tabsSelected = 0
// 			createContent(tabsSelected, listSelected)
// 			go table(contentPos)
// 		case 'f':
// 			if tabsSelected == 1 {
// 				continue
// 			}
// 			tabsSelected = 1
// 			createContent(tabsSelected, listSelected)
// 			go matchday(contentPos, selectedMatchday)
// 		case 's':
// 			if tabsSelected == 2 {
// 				continue
// 			}
// 			tabsSelected = 2
// 			createContent(tabsSelected, listSelected)
// 		case 'p':
// 			if tabsSelected == 3 {
// 				continue
// 			}
// 			tabsSelected = 3
// 			createContent(tabsSelected, listSelected)
// 		case 'k':
// 			if colored && selectedMatchday-1 >= 0 {
// 				selectedMatchday--
// 				go matchday(4, 20, selectedMatchday)
// 				break
// 			}
// 			if listSelected-1 >= 0 {
// 				listSelected--
// 				createContent(tabsSelected, listSelected)
// 			}
// 		case 'j':
// 			if colored && selectedMatchday+1 <= 34 {
// 				selectedMatchday++
// 				go matchday(4, 20, selectedMatchday)
// 				break
// 			}
// 			if listSelected+1 <= 5 {
// 				listSelected++
// 				createContent(tabsSelected, listSelected)
// 			}
// 		case '\r':
// 			colored = true
// 			createContent(tabsSelected, listSelected)
// 			go matchday(4, 20, selectedMatchday)
// 		case 'B':
// 		case 'A':
// 		case 27:
// 			colored = false
// 			createContent(tabsSelected, listSelected)
// 		}
// 	}
//
// 	component.TearDown()
// }
//
// func firstVersion() {
// 	component.Start()
// 	defer component.TearDown()
// 	fd := int(syscall.Stdin)
//
// 	oldState, err := term.MakeRaw(fd)
// 	if err != nil {
// 		fmt.Println("Error setting raw mode:", err)
// 		return
// 	}
// 	defer term.Restore(fd, oldState)
//
// 	boxMatchday := component.NewBox("Matchday").WithPadding().WithRoundedCorners().WithColoredBorder(component.Orange).At(5, 40)
// 	boxTable := component.NewBox("Table").WithPadding().WithRoundedCorners().At(5, 55)
// 	component.Print(boxMatchday)
// 	component.Print(boxTable)
// 	go matchday(10, 20, 2)
//
// 	selected := 0
//
// 	reader := bufio.NewReader(os.Stdin)
//
// outer:
//
// 	for {
// 		b, err := reader.ReadByte()
// 		if err != nil {
// 			fmt.Println("Error reading byte:", err)
// 			break
// 		}
//
// 		switch b {
// 		case 'q', component.CtrlC:
// 			break outer
// 		case 'm':
// 			if selected == 0 {
// 				break
// 			}
// 			selected = 0
// 			boxMatchday = component.NewBox("Matchday").WithPadding().WithRoundedCorners().WithColoredBorder(orange)
// 			boxTable = component.NewBox("Table").WithPadding().WithRoundedCorners()
// 			go matchday(component.NewPos(10, 20), 2)
// 		case 't':
// 			if selected == 1 {
// 				break
// 			}
// 			selected = 1
// 			boxMatchday = component.NewBox("Matchday").WithPadding().WithRoundedCorners()
// 			boxTable = component.NewBox("Table").WithPadding().WithRoundedCorners().WithColoredBorder(component.orange)
// 			go table(component.NewPos(10, 20))
// 		}
// 		component.Clear()
// 		component.Print(component.NewPos(5, 40), boxMatchday.String())
// 		component.Print(NewPos(5, 55), boxTable.String())
// 	}
// }
//
// func matchday(row, col, matchNum int) {
// 	matches := getMatchday("bl1", matchNum)
// 	matchStrs := make([][]string, 0, 9)
// 	maxLenHome := 0
// 	for _, match := range matches {
// 		var (
// 			finished  = match.MatchIsFinished
// 			home      = match.Home.TeamName
// 			away      = match.Away.TeamName
// 			result    MatchResults
// 			homeGoals int
// 			awayGoals int
// 		)
//
// 		if finished {
// 			result = match.MatchResults[1]
// 			homeGoals = result.PointsHome
// 			awayGoals = result.PointsAway
// 		}
//
// 		if utf8.RuneCountInString(home) > maxLenHome {
// 			maxLenHome = utf8.RuneCountInString(home)
// 		}
//
// 		if finished {
// 			matchStrs = append(matchStrs, []string{home, strconv.Itoa(homeGoals), strconv.Itoa(awayGoals), away})
// 		} else {
// 			matchStrs = append(matchStrs, []string{home, "-", "-", away})
// 		}
// 	}
//
// 	content := make([]string, 9, 9)
// 	for i, c := range matchStrs {
// 		var (
// 			home      = c[0]
// 			homeGoals = c[1]
// 			awayGoals = c[2]
// 			away      = c[3]
// 		)
// 		b := new(strings.Builder)
//
// 		b.WriteString(strings.Repeat(" ", maxLenHome-utf8.RuneCountInString(home)))
// 		b.WriteString(fmt.Sprintf("%s %s : %s %s", home, homeGoals, awayGoals, away))
//
// 		content[i] = b.String()
// 	}
//
// 	box := NewBox(content...).WithRoundedCorners().WithTitle(fmt.Sprintf("%s %d", "Matchday", matchNum)).WithPadding(1)
// 	Print(row, col, box.String())
// }
//
// func table() {
// 	tableInfo := getTable()
// 	table := NewTable(
// 		NewHeader("#", true),
// 		NewHeader("TeamName", false),
// 		NewHeader("Games", true),
// 		NewHeader("W", true),
// 		NewHeader("D", true),
// 		NewHeader("L", true),
// 		NewHeader("Goals", true),
// 		NewHeader("Diff", true),
// 		NewHeader("Points", true),
// 	).WithRoundedCorners()
//
// 	for i, info := range tableInfo {
// 		data := fmt.Sprintf("%d??%s??%d??%d??%d??%d??%d??%d??%d", i+1, info.TeamName, info.Matches, info.Won, info.Draw, info.Lost, info.Goals, info.GoalDiff, info.Points)
// 		table.AddRow(strings.Split(data, "??"))
// 	}
//
// 	Print(table)
// }
