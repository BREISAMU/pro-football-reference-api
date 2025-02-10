package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PlayoffGame struct {
	Round      int    `json:"roundNum"`
	Week       string `json:"roundName"`
	Day        string `json:"day"`
	Date       string `json:"date"`
	Winner     string `json:"winner"`
	Loser      string `json:"loser"`
	BoxScoreId string `json:"boxScoreId"`
	PointsW    int    `json:"ptsW"`
	PointsL    int    `json:"ptsL"`
}

type PlayoffResult struct {
	Year  int           `json:"year"`
	Champ string        `json:"champ"`
	Games []PlayoffGame `json:"games"`
}

func GetPlayoffResultsByYear(url string, year int) (PlayoffResult, error) {
	tableSelector := "#playoff_results"

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
			return PlayoffResult{}, fmt.Errorf("error creating request: %v", err)
		}

		// Headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return PlayoffResult{}, fmt.Errorf("error making request: %v", err)
		}

		// Rate limit check
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return PlayoffResult{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
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
		return PlayoffResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return PlayoffResult{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	var tableData [][]string
	results := PlayoffResult{year, "Not found yet", []PlayoffGame{}}

	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})

		tableData = append(tableData, rowData)
		if err == nil {
			round := 1
			week := rowData[0]
			day := rowData[1]
			date := rowData[2]
			winner := rowData[3]
			loser := rowData[4]
			boxScoreId := rowData[5]
			ptsW, _ := strconv.Atoi(rowData[6])
			ptsL, _ := strconv.Atoi(rowData[7])

			game := PlayoffGame{
				Round:      round,
				Week:       week,
				Day:        day,
				Date:       date,
				Winner:     winner,
				Loser:      loser,
				BoxScoreId: boxScoreId,
				PointsW:    ptsW,
				PointsL:    ptsL,
			}

			results.Games = append(results.Games, game)
		}
	})

	if len(tableData) < 1 {
		return PlayoffResult{}, fmt.Errorf("no data found for selected year")
	}

	return results, nil
}
