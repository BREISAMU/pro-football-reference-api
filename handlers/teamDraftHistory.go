package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type DraftPick struct {
	Year          int    `json:"year"`
	Round         int    `json:"round"`
	Name          string `json:"name"`
	Pick          int    `json:"pick"`
	Position      string `json:"position"`
	LastSeason    int    `json:"lastSeason"`
	FirstAllPro   int    `json:"firstAllPro"`
	ProBowl       int    `json:"proBowl"`
	StarterYears  int    `json:"starterYears"`
	CareerAV      int    `json:"careerAv"`
	GamesPlayed   int    `json:"gamesPlayed"`
	PassCmp       int    `json:"passCmp"`
	PassAtt       int    `json:"passAtt"`
	PassYds       int    `json:"passYds"`
	PassTDs       int    `json:"passTds"`
	PassInts      int    `json:"passInts"`
	RushAtt       int    `json:"rushAtt"`
	RushYds       int    `json:"rushYds"`
	RushTDs       int    `json:"rushTds"`
	ReceivingRecs int    `json:"receivingRecs"`
	ReceivingYds  int    `json:"receivingYds"`
	ReceivingTDs  int    `json:"receivingTds"`
	DefInts       int    `json:"defInts"`
	DefSacks      int    `json:"defSacks"`
	College       string `json:"college"`
}

func GetDraftYear(url string, tableSelector string, year int) ([]DraftPick, error) {
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
			return []DraftPick{}, fmt.Errorf("error creating request: %v", err)
		}

		// Headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")

		resp, err = client.Do(req)
		if err != nil {
			return []DraftPick{}, fmt.Errorf("error making request: %v", err)
		}

		// Rate limit check
		if resp.StatusCode == 429 {
			resp.Body.Close()
			if attempt == maxRetries {
				return []DraftPick{}, fmt.Errorf("hit rate limit after %d attempts", maxRetries)
			}

			retryAfter := resp.Header.Get("Retry-After")
			waitTime := 60 * time.Second
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
		return []DraftPick{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return []DraftPick{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	var tableData [][]string

	// ---- END ----

	var draft [][]string
	doc.Find(tableSelector).Find("tr").Each(func(i int, row *goquery.Selection) {
		var rowData []string
		row.Find("td, th").Each(func(j int, cell *goquery.Selection) {
			rowData = append(rowData, cell.Text())
		})
		if len(rowData) > 0 {
			tableData = append(tableData, rowData)
			rowYear, err := strconv.Atoi(rowData[0])
			if err == nil && rowYear == year {
				draft = append(draft, rowData)
			}
		}
	})

	if len(draft) == 0 {
		return []DraftPick{}, fmt.Errorf("no data found for year %d", year)
	}

	resDraft := []DraftPick{}

	for i := 0; i < len(draft); i++ {
		fmt.Println(draft[i])
		year, _ := strconv.Atoi(draft[i][0])
		round, _ := strconv.Atoi(draft[i][1])
		name := draft[i][2]
		pick, _ := strconv.Atoi(draft[i][3])
		position := draft[i][4]
		lastSeason, _ := strconv.Atoi(draft[i][5])
		firstAllPro, _ := strconv.Atoi(draft[i][6])
		proBowl, _ := strconv.Atoi(draft[i][7])
		starterYears, _ := strconv.Atoi(draft[i][8])
		careerAv, _ := strconv.Atoi(draft[i][9])
		gamesPlayed, _ := strconv.Atoi(draft[i][10])
		passCmp, _ := strconv.Atoi(draft[i][11])
		passAtt, _ := strconv.Atoi(draft[i][12])
		passYds, _ := strconv.Atoi(draft[i][13])
		passTds, _ := strconv.Atoi(draft[i][14])
		passInts, _ := strconv.Atoi(draft[i][15])
		rushAtt, _ := strconv.Atoi(draft[i][16])
		rushYds, _ := strconv.Atoi(draft[i][17])
		rushTds, _ := strconv.Atoi(draft[i][18])
		receivingRecs, _ := strconv.Atoi(draft[i][19])
		receivingYds, _ := strconv.Atoi(draft[i][20])
		receivingTds, _ := strconv.Atoi(draft[i][21])
		defInts, _ := strconv.Atoi(draft[i][22])
		defSks, _ := strconv.Atoi(draft[i][23])
		college := draft[i][24]

		draftPick := DraftPick{
			Year:          year,
			Round:         round,
			Name:          name,
			Pick:          pick,
			Position:      position,
			LastSeason:    lastSeason,
			FirstAllPro:   firstAllPro,
			ProBowl:       proBowl,
			StarterYears:  starterYears,
			CareerAV:      careerAv,
			GamesPlayed:   gamesPlayed,
			PassCmp:       passCmp,
			PassAtt:       passAtt,
			PassYds:       passYds,
			PassTDs:       passTds,
			PassInts:      passInts,
			RushAtt:       rushAtt,
			RushYds:       rushYds,
			RushTDs:       rushTds,
			ReceivingRecs: receivingRecs,
			ReceivingYds:  receivingYds,
			ReceivingTDs:  receivingTds,
			DefInts:       defInts,
			DefSacks:      defSks,
			College:       college,
		}

		resDraft = append(resDraft, draftPick)
	}

	return resDraft, nil
}
