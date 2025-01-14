package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Stats struct {
	DataType                string  `json:"dataType"`
	PointsFor               int     `json:"pointsFor"`
	TotalYards              int     `json:"totalYards"`
	TotalPlays              int     `json:"totalPlays"`
	YardsPerPlay            float64 `json:"yardsPerPlay"`
	Turnovers               int     `json:"turnovers"`
	Fumbles                 int     `json:"fumbles"`
	FirstDowns              int     `json:"firstDowns"`
	PassCompletions         int     `json:"passCompletions"`
	PassAttempts            int     `json:"passAttempts"`
	PassYards               int     `json:"passYards"`
	PassTds                 int     `json:"passTds"`
	PassInts                int     `json:"passInts"`
	PassYardsPerAtt         float64 `json:"passYardsPerAtt"`
	PassFirstDowns          int     `json:"passFirstDowns"`
	RushAttempts            int     `json:"rushAttempts"`
	RushYards               int     `json:"rushYards"`
	RushTDs                 int     `json:"rushTDs"`
	RushYardsPerAtt         float64 `json:"rushYardsPerAtt"`
	RushFirstDowns          int     `json:"rushFirstDowns"`
	Penalties               int     `json:"penalties"`
	PenaltyYards            int     `json:"penaltyYards"`
	PenaltyFirstDowns       int     `json:"penaltyFirstDowns"`
	Drives                  int     `json:"drives"`
	ScoringDrivePercentage  float64 `json:"scoringDrivePercentage"`
	TurnoverDrivePercentage float64 `json:"turnoverDrivePercentage"`
	AverageStartPosition    float64 `json:"averageStartPosition"`
	AvgDriveLength          float64 `json:"avgDriveLength"`
	AvgDrivePlays           float64 `json:"avgDrivePlays"`
	AvgDriveYards           float64 `json:"avgDriveYards"`
	AvgDrivePoints          float64 `json:"avgDrivePoints"`
}

// int only rankings
type Rankings struct {
	DataType                string `json:"dataType"`
	PointsFor               int    `json:"pointsFor"`
	TotalYards              int    `json:"totalYards"`
	TotalPlays              int    `json:"totalPlays"`
	YardsPerPlay            int    `json:"yardsPerPlay"`
	Turnovers               int    `json:"turnovers"`
	Fumbles                 int    `json:"fumbles"`
	FirstDowns              int    `json:"firstDowns"`
	PassCompletions         int    `json:"passCompletions"`
	PassAttempts            int    `json:"passAttempts"`
	PassYards               int    `json:"passYards"`
	PassTds                 int    `json:"passTds"`
	PassInts                int    `json:"passInts"`
	PassYardsPerAtt         int    `json:"passYardsPerAtt"`
	PassFirstDowns          int    `json:"passFirstDowns"`
	RushAttempts            int    `json:"rushAttempts"`
	RushYards               int    `json:"rushYards"`
	RushTDs                 int    `json:"rushTDs"`
	RushYardsPerAtt         int    `json:"rushYardsPerAtt"`
	RushFirstDowns          int    `json:"rushFirstDowns"`
	Penalties               int    `json:"penalties"`
	PenaltyYards            int    `json:"penaltyYards"`
	PenaltyFirstDowns       int    `json:"penaltyFirstDowns"`
	Drives                  int    `json:"drives"`
	ScoringDrivePercentage  int    `json:"scoringDrivePercentage"`
	TurnoverDrivePercentage int    `json:"turnoverDrivePercentage"`
	AverageStartPosition    int    `json:"averageStartPosition"`
	AvgDriveLength          int    `json:"avgDriveLength"`
	AvgDrivePlays           int    `json:"avgDrivePlays"`
	AvgDriveYards           int    `json:"avgDriveYards"`
	AvgDrivePoints          int    `json:"avgDrivePoints"`
}

func GetTeamYearStats(url string, tableSelector string) (Stats, Stats, Rankings, Rankings, error) {
	// ---- CLIENT BOILERPLATE ----
	client := &http.Client{
		Timeout: 4 * (time.Second + 8),
	}

	maxRetries := 2
	var resp *http.Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("error creating request: %v", err)
		}

		// Headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("error making request: %v", err)
		}

		// Rate limit check
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 15 * time.Second
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					waitTime = time.Duration(seconds) * time.Second
				}
			}

			log.Printf("Rate limited. Waiting %v before retry %d/%d", waitTime, attempt, maxRetries)
			time.Sleep(waitTime)
			continue
		}

		// Successful response
		if resp.StatusCode == 200 {
			break
		}

		resp.Body.Close()
		return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	var tableData [][]string

	// ---- END ----

	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})
		if len(rowData) > 0 {
			tableData = append(tableData, rowData)
		}
	})

	if len(tableData) < 6 {
		return Stats{}, Stats{}, Rankings{}, Rankings{}, fmt.Errorf("no data found for selected year")
	}

	resStats := []Stats{}
	resRankings := []Rankings{}
	for i := 2; i < 4; i++ {
		statGroup := tableData[i]

		dataType := statGroup[0]
		pointsFor, _ := strconv.Atoi(statGroup[1])
		totalYards, _ := strconv.Atoi(statGroup[2])
		totalPlays, _ := strconv.Atoi(statGroup[3])
		yardsPerPlay, _ := strconv.ParseFloat(statGroup[4], 64)
		turnovers, _ := strconv.Atoi(statGroup[5])
		fumbles, _ := strconv.Atoi(statGroup[6])
		firstDowns, _ := strconv.Atoi(statGroup[7])
		passCompletions, _ := strconv.Atoi(statGroup[8])
		passAttempts, _ := strconv.Atoi(statGroup[9])
		passYards, _ := strconv.Atoi(statGroup[10])
		passTds, _ := strconv.Atoi(statGroup[11])
		passInts, _ := strconv.Atoi(statGroup[12])
		passYardsPerAtt, _ := strconv.ParseFloat(statGroup[13], 64)
		passFirstDowns, _ := strconv.Atoi(statGroup[14])
		rushAttempts, _ := strconv.Atoi(statGroup[15])
		rushYards, _ := strconv.Atoi(statGroup[16])
		rushTDs, _ := strconv.Atoi(statGroup[17])
		rushYardsPerAtt, _ := strconv.ParseFloat(statGroup[18], 64)
		rushFirstDowns, _ := strconv.Atoi(statGroup[19])
		penalties, _ := strconv.Atoi(statGroup[20])
		penaltyYards, _ := strconv.Atoi(statGroup[21])
		penaltyFirstDowns, _ := strconv.Atoi(statGroup[22])
		drives, _ := strconv.Atoi(statGroup[23])
		scoringDrivePercentage, _ := strconv.ParseFloat(statGroup[24], 64)
		turnoverDrivePercentage, _ := strconv.ParseFloat(statGroup[25], 64)

		averageStartPosition, _ := strconv.ParseFloat(statGroup[26][len(statGroup[26])-4:], 64)

		hour, _ := strconv.ParseFloat(string(statGroup[27][0]), 64)
		minute, _ := strconv.ParseFloat(statGroup[27][2:], 64)
		avgDriveLength := hour + (minute / 60)

		avgDrivePlays, _ := strconv.ParseFloat(statGroup[28], 64)
		avgDriveYards, _ := strconv.ParseFloat(statGroup[29], 64)
		avgDrivePoints, _ := strconv.ParseFloat(statGroup[30], 64)

		res := Stats{
			DataType:                dataType,
			PointsFor:               pointsFor,
			TotalYards:              totalYards,
			TotalPlays:              totalPlays,
			YardsPerPlay:            yardsPerPlay,
			Turnovers:               turnovers,
			Fumbles:                 fumbles,
			FirstDowns:              firstDowns,
			PassCompletions:         passCompletions,
			PassAttempts:            passAttempts,
			PassYards:               passYards,
			PassTds:                 passTds,
			PassInts:                passInts,
			PassYardsPerAtt:         passYardsPerAtt,
			PassFirstDowns:          passFirstDowns,
			RushAttempts:            rushAttempts,
			RushYards:               rushYards,
			RushTDs:                 rushTDs,
			RushYardsPerAtt:         rushYardsPerAtt,
			RushFirstDowns:          rushFirstDowns,
			Penalties:               penalties,
			PenaltyYards:            penaltyYards,
			PenaltyFirstDowns:       penaltyFirstDowns,
			Drives:                  drives,
			ScoringDrivePercentage:  scoringDrivePercentage,
			TurnoverDrivePercentage: turnoverDrivePercentage,
			AverageStartPosition:    averageStartPosition,
			AvgDriveLength:          avgDriveLength,
			AvgDrivePlays:           avgDrivePlays,
			AvgDriveYards:           avgDriveYards,
			AvgDrivePoints:          avgDrivePoints,
		}

		resStats = append(resStats, res)
	}

	for i := 4; i < 6; i++ {
		statGroup := tableData[i]

		dataType := statGroup[0]
		pointsFor, _ := strconv.Atoi(statGroup[1])
		totalYards, _ := strconv.Atoi(statGroup[2])
		totalPlays, _ := strconv.Atoi(statGroup[3])
		yardsPerPlay, _ := strconv.Atoi(statGroup[4])
		turnovers, _ := strconv.Atoi(statGroup[5])
		fumbles, _ := strconv.Atoi(statGroup[6])
		firstDowns, _ := strconv.Atoi(statGroup[7])
		passCompletions, _ := strconv.Atoi(statGroup[8])
		passAttempts, _ := strconv.Atoi(statGroup[9])
		passYards, _ := strconv.Atoi(statGroup[10])
		passTds, _ := strconv.Atoi(statGroup[11])
		passInts, _ := strconv.Atoi(statGroup[12])
		passYardsPerAtt, _ := strconv.Atoi(statGroup[13])
		passFirstDowns, _ := strconv.Atoi(statGroup[14])
		rushAttempts, _ := strconv.Atoi(statGroup[15])
		rushYards, _ := strconv.Atoi(statGroup[16])
		rushTDs, _ := strconv.Atoi(statGroup[17])
		rushYardsPerAtt, _ := strconv.Atoi(statGroup[18])
		rushFirstDowns, _ := strconv.Atoi(statGroup[19])
		penalties, _ := strconv.Atoi(statGroup[20])
		penaltyYards, _ := strconv.Atoi(statGroup[21])
		penaltyFirstDowns, _ := strconv.Atoi(statGroup[22])
		drives, _ := strconv.Atoi(statGroup[23])
		scoringDrivePercentage, _ := strconv.Atoi(statGroup[24])
		turnoverDrivePercentage, _ := strconv.Atoi(statGroup[25])
		averageStartPosition, _ := strconv.Atoi(statGroup[26])
		avgDriveLength, _ := strconv.Atoi(statGroup[27])
		avgDrivePlays, _ := strconv.Atoi(statGroup[28])
		avgDriveYards, _ := strconv.Atoi(statGroup[29])
		avgDrivePoints, _ := strconv.Atoi(statGroup[30])

		res := Rankings{
			DataType:                dataType,
			PointsFor:               pointsFor,
			TotalYards:              totalYards,
			TotalPlays:              totalPlays,
			YardsPerPlay:            yardsPerPlay,
			Turnovers:               turnovers,
			Fumbles:                 fumbles,
			FirstDowns:              firstDowns,
			PassCompletions:         passCompletions,
			PassAttempts:            passAttempts,
			PassYards:               passYards,
			PassTds:                 passTds,
			PassInts:                passInts,
			PassYardsPerAtt:         passYardsPerAtt,
			PassFirstDowns:          passFirstDowns,
			RushAttempts:            rushAttempts,
			RushYards:               rushYards,
			RushTDs:                 rushTDs,
			RushYardsPerAtt:         rushYardsPerAtt,
			RushFirstDowns:          rushFirstDowns,
			Penalties:               penalties,
			PenaltyYards:            penaltyYards,
			PenaltyFirstDowns:       penaltyFirstDowns,
			Drives:                  drives,
			ScoringDrivePercentage:  scoringDrivePercentage,
			TurnoverDrivePercentage: turnoverDrivePercentage,
			AverageStartPosition:    averageStartPosition,
			AvgDriveLength:          avgDriveLength,
			AvgDrivePlays:           avgDrivePlays,
			AvgDriveYards:           avgDriveYards,
			AvgDrivePoints:          avgDrivePoints,
		}

		resRankings = append(resRankings, res)
	}

	return resStats[0], resStats[1], resRankings[0], resRankings[1], nil
}
