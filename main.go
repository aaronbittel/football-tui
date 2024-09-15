package main

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
