package main

import "time"

type MatchAPI struct {
	MatchID            int            `json:"matchID"`
	MatchDateTime      string         `json:"matchDateTime"`
	TimeZoneID         string         `json:"timeZoneID"`
	LeagueID           int            `json:"leagueId"`
	LeagueName         string         `json:"leagueName"`
	LeagueSeason       int            `json:"leagueSeason"`
	LeagueShortcut     string         `json:"leagueShortcut"`
	MatchDateTimeUTC   time.Time      `json:"matchDateTimeUTC"`
	Group              Group          `json:"group"`
	Home               Team           `json:"team1"`
	Away               Team           `json:"team2"`
	LastUpdateDateTime string         `json:"lastUpdateDateTime"`
	MatchIsFinished    bool           `json:"matchIsFinished"`
	MatchResults       []MatchResults `json:"matchResults"`
	Goals              []Goals        `json:"goals"`
	Location           Location       `json:"location"`
	NumberOfViewers    any            `json:"numberOfViewers"`
}

type Group struct {
	GroupName    string `json:"groupName"`
	GroupOrderID int    `json:"groupOrderID"`
	GroupID      int    `json:"groupID"`
}

type Team struct {
	TeamID        int    `json:"teamId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconUrl"`
	TeamGroupName any    `json:"teamGroupName"`
}

type MatchResults struct {
	ResultID          int    `json:"resultID"`
	ResultName        string `json:"resultName"`
	PointsHome        int    `json:"pointsTeam1"`
	PointsAway        int    `json:"pointsTeam2"`
	ResultOrderID     int    `json:"resultOrderID"`
	ResultTypeID      int    `json:"resultTypeID"`
	ResultDescription string `json:"resultDescription"`
}

type Goals struct {
	GoalID         int    `json:"goalID"`
	ScoreHome      int    `json:"scoreTeam1"`
	ScoreAway      int    `json:"scoreTeam2"`
	MatchMinute    int    `json:"matchMinute"`
	GoalGetterID   int    `json:"goalGetterID"`
	GoalGetterName string `json:"goalGetterName"`
	IsPenalty      bool   `json:"isPenalty"`
	IsOwnGoal      bool   `json:"isOwnGoal"`
	IsOvertime     bool   `json:"isOvertime"`
	Comment        any    `json:"comment"`
}

type Location struct {
	LocationID      int    `json:"locationID"`
	LocationCity    string `json:"locationCity"`
	LocationStadium string `json:"locationStadium"`
}
