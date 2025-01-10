package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type SeasonOverlook struct {
	Year               int
	League             string
	Team               string
	Wins               int
	Losses             int
	Ties               int
	DivisionFinish     int
	PlayoffExitRound   int
	PointsFor          int
	PointsAgainst      int
	pointsDif          int
	HeadCoaches        string
	BestPlayerAv       string
	BestPlayerPasser   string
	BestPlayerRusher   string
	BestPlayerReceiver string
	OffRankPts         int
	OffRankYds         int
	DefRankPts         int
	DefRankYds         int
	TakeawayRank       int
	PointsDifRank      int
	YardsDifRank       int
	TeamsInLeague      int
	MarginOfVictory    float64
	StrengthOfSchedule float64
	Srs                float64
	OffensiveSrs       float64
	DefensiveSrs       float64
}

func main() {
	team := "gnb"
	url := "https://www.pro-football-reference.com/teams/" + team + "/"
	tableSelector := "#team_index"

	seasonA, errA := ScrapeTeamSeasonOverlook(url, tableSelector, 1980)
	if errA != nil {
		log.Println(errA)
	} else {
		fmt.Println(seasonA)
	}
}

func ScrapeTeamSeasonOverlook(url string, tableSelector string, year int) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "{}", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "{}", err
	}

	var tableData [][]string
	var season []string
	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})
		tableData = append(tableData, rowData)

		rowYear, _ := strconv.Atoi(rowData[0])

		if rowYear == year {
			season = rowData
		}
	})

	if len(season) == 0 {
		return "{}", err
	}

	league := season[1]
	team := season[2]
	wins, _ := strconv.Atoi(season[3])
	losses, _ := strconv.Atoi(season[4])
	ties, _ := strconv.Atoi(season[5])
	divisionFinish, _ := strconv.Atoi(season[6][0:1])

	playoffExitRoundString := season[7]
	playoffExitRoundInt := 0
	if playoffExitRoundString == "Lost WC" {
		playoffExitRoundInt = 1
	} else if playoffExitRoundString == "Lost Div" {
		playoffExitRoundInt = 2
	} else if playoffExitRoundString == "Lost Conf" {
		playoffExitRoundInt = 3
	} else if playoffExitRoundString == "Lost SB" || playoffExitRoundString == "Lost Champ" {
		playoffExitRoundInt = 4
	} else if playoffExitRoundString == "Won SB" || playoffExitRoundString == "Won Champ" {
		playoffExitRoundInt = 5
	}

	pointsFor, _ := strconv.Atoi(season[8])
	pointsAgainst, _ := strconv.Atoi(season[9])
	pointsDif, _ := strconv.Atoi(season[10])
	headCoaches := season[11]
	bestPlayerAv := season[12]
	bestPlayerPasser := season[13]
	bestPlayerRusher := season[14]
	bestPlayerReceiver := season[15]
	offRankPts, _ := strconv.Atoi(season[16])
	offRankYds, _ := strconv.Atoi(season[17])
	defRankPts, _ := strconv.Atoi(season[18])
	defRankYds, _ := strconv.Atoi(season[19])
	takeawayRank, _ := strconv.Atoi(season[20])
	pointsDifRank, _ := strconv.Atoi(season[21])
	yardsDifRank, _ := strconv.Atoi(season[22])
	teamsInLeague, _ := strconv.Atoi(season[23])
	marginOfVictory, _ := strconv.ParseFloat(season[24], 64)
	strengthOfSchedule, _ := strconv.ParseFloat(season[25], 64)
	srs, _ := strconv.ParseFloat(season[26], 64)
	offensiveSrs, _ := strconv.ParseFloat(season[27], 64)
	defensiveSrs, _ := strconv.ParseFloat(season[28], 64)

	seasonOverlook := SeasonOverlook{
		Year:               year,
		League:             league,
		Team:               team,
		Wins:               wins,
		Losses:             losses,
		Ties:               ties,
		DivisionFinish:     divisionFinish,
		PlayoffExitRound:   playoffExitRoundInt,
		PointsFor:          pointsFor,
		PointsAgainst:      pointsAgainst,
		pointsDif:          pointsDif,
		HeadCoaches:        headCoaches,
		BestPlayerAv:       bestPlayerAv,
		BestPlayerPasser:   bestPlayerPasser,
		BestPlayerRusher:   bestPlayerRusher,
		BestPlayerReceiver: bestPlayerReceiver,
		OffRankPts:         offRankPts,
		OffRankYds:         offRankYds,
		DefRankPts:         defRankPts,
		DefRankYds:         defRankYds,
		TakeawayRank:       takeawayRank,
		PointsDifRank:      pointsDifRank,
		YardsDifRank:       yardsDifRank,
		TeamsInLeague:      teamsInLeague,
		MarginOfVictory:    marginOfVictory,
		StrengthOfSchedule: strengthOfSchedule,
		Srs:                srs,
		OffensiveSrs:       offensiveSrs,
		DefensiveSrs:       defensiveSrs,
	}

	seasonOverlookJson, err := json.Marshal(seasonOverlook)
	if err != nil {
		fmt.Println(err)
		return "{}", nil
	}

	return string(seasonOverlookJson), nil
}
