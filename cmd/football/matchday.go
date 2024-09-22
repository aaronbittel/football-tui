package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// )
//
// type Match struct {
// 	Home      string
// 	Away      string
// 	HomeGoals int
// 	AwayGoals int
// 	Finished  bool
// }
//
// func NewMatch(home, away string, homeGoals, awayGoals int, finished bool) Match {
// 	return Match{
// 		Home:      home,
// 		Away:      away,
// 		HomeGoals: homeGoals,
// 		AwayGoals: awayGoals,
// 		Finished:  finished,
// 	}
// }
//
// func (m Match) String() string {
// 	if m.Finished {
// 		return fmt.Sprintf("%s %d : %d %s", m.Home, m.HomeGoals, m.AwayGoals, m.Away)
// 	}
// 	return fmt.Sprintf("%s - : - %s", m.Home, m.Away)
// }
//
// func FromMatchApi(matchesApi []MatchAPI) []Match {
// 	var (
// 		matches = []Match{}
// 		m       Match
// 	)
// 	for _, mApi := range matchesApi {
// 		var (
// 			home      = mApi.Home.TeamName
// 			away      = mApi.Away.TeamName
// 			homeGoals int
// 			awayGoals int
// 			finished  = len(mApi.MatchResults) > 0
// 		)
// 		if len(mApi.MatchResults) > 0 {
// 			results := mApi.MatchResults[1]
// 			homeGoals = results.PointsHome
// 			awayGoals = results.PointsAway
// 		}
//
// 		m = NewMatch(home, away, homeGoals, awayGoals, finished)
// 		matches = append(matches, m)
// 	}
// 	return matches
// }
//
// func getMatchday(league string, i int) []Match {
// 	matchday := make([]Match, 0, 9)
//
// 	url := "https://api.openligadb.de/getmatchdata/%s/2024/%d"
// 	resp, err := http.Get(fmt.Sprintf(url, league, i))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
//
// 	decoder := json.NewDecoder(resp.Body)
// 	err = decoder.Decode(&matchday)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	return matchday
// }
//
// func getMatchdayFromFile(i int) []Match {
// 	matches := make([]MatchAPI, 34*9)
//
// 	f, err := os.Open("data.json")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	decoder := json.NewDecoder(f)
// 	decoder.Decode(&matches)
// 	return FromMatchApi(matches[(i-1)*9 : i*9])
// }
