package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type TableInfo struct {
	TeamInfoID    int    `json:"teamInfoId"`
	TeamName      string `json:"teamName"`
	ShortName     string `json:"shortName"`
	TeamIconURL   string `json:"teamIconUrl"`
	Points        int    `json:"points"`
	OpponentGoals int    `json:"opponentGoals"`
	Goals         int    `json:"goals"`
	Matches       int    `json:"matches"`
	Won           int    `json:"won"`
	Lost          int    `json:"lost"`
	Draw          int    `json:"draw"`
	GoalDiff      int    `json:"goalDiff"`
}

func getTable() []TableInfo {
	table := make([]TableInfo, 0, 9)

	resp, err := http.Get("https://api.openligadb.de/getbltable/bl1/2024")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&table)
	if err != nil {
		log.Fatal(err)
	}

	return table
}
